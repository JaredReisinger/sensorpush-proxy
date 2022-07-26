package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
)

const (
	appName = "sensorpush-proxy"
)

type config struct {
	SensorPush struct {
		Username  string
		Password  string
		DeviceIDs map[string]string
	}
	Proxy struct {
		Port string
	}
}

func main() {
	log.Print(appName)

	// viper.SetDefault("sensorPush.username", "")
	// viper.SetDefault("sensorPush.password", "")
	// viper.SetDefault("sensorPush.deviceId", "")
	viper.SetDefault("proxy.port", ":5375")

	viper.SetConfigName("config")
	// viper.SetConfigType("yaml") // setting the config type takes precedence
	// over the extension, which seems wrong!

	viper.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("SENSORPUSH")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		case *viper.ConfigFileNotFoundError:
			log.Printf("no config file found... perhaps there are environment variables")
		default:
			log.Fatalf("ERROR: (%T) %s", err, err.Error())
		}
	}

	var c config
	viper.Unmarshal(&c)
	log.Printf("unmarshaled config: %+v", c)

	invalidArgs := []string{}

	user := viper.GetString("sensorPush.username")
	if user == "" {
		invalidArgs = append(invalidArgs, "--username (SENSORPUSH_USERNAME)")
	}

	pass := viper.GetString("sensorPush.password")
	if pass == "" {
		invalidArgs = append(invalidArgs, "--password (SENSORPUSH_PASSWORD)")
	}

	// deviceID := viper.GetString("sensorPush.deviceIds")
	port := viper.GetString("proxy.port")
	if port == "" {
		invalidArgs = append(invalidArgs, "--port (SENSORPUSH_PORT)")
	}

	if len(invalidArgs) > 0 {
		log.Fatalf("invalid or missing args:\n  %s", strings.Join(invalidArgs, "\n  "))
	}

	// if user == "" || pass == "" || /* deviceID == "" || */ port == "" {
	// 	log.Fatalf("one of SENSORPUSH_USERNAME (%q), SENSORPUSH_PASSWORD (length %d), SENSORPUSH_DEVICE_ID (%q), or SENSORPUSH_PORT (%q) is missing", user, len(pass), deviceID, port)
	// }

	// ensure port has a ":" prefix?...
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	client, err := sensorpush.NewClient(user, pass)
	if err != nil {
		log.Fatalf("unable to create client: %+v", err)
	}

	// TODO: should lastSample/SuccessfulCall really be channels?  We don't
	// really want to server/read from them *while* updates are happening...
	var lastSample sensorpush.Sample
	lastSuccessfulCall := time.Now()

	appCtx, appCancel := context.WithCancel(context.Background())

	updater := func() {
		sample, err := client.LastSample("deviceID")
		if err != nil {
			log.Printf("unable to get sample: %+v", err)
			// if we're X? past the last successful call, it's time to cancel
			// and exit, and let a new process/container start up
			if time.Since(lastSuccessfulCall) > (5 * time.Minute) {
				appCancel()
			}
		} else {
			log.Printf(
				"UPDATER: %s: %.1f°F, %.1f%%RH",
				sample.Observed.Local().Format(time.RFC3339),
				sample.Temperature,
				sample.Humidity,
			)
			lastSample = *sample
			lastSuccessfulCall = time.Now()
		}
	}

	runBackground(appCtx, updater, time.Minute)

	// Now spin up a web server to serve the sample data...
	http.HandleFunc("/sensor", func(w http.ResponseWriter, req *http.Request) {
		// If there's an Origin header, send back Access-Control-Allow-Origin
		origin := req.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		b, err := json.Marshal(lastSample)
		if err != nil {
			io.WriteString(w, "{ error: \"no data?\" }")
			return
		}
		w.Write(b)
	})

	srv := &http.Server{
		Addr: port,
		// Handler: ... http.DefaultServeMux, which uses the above handler
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	// also run a one-off background goroutine that closes the server when the
	// context is canceled...
	shutdownContext, shutdownComplete := context.WithCancel(context.Background())
	go func() {
		<-appCtx.Done()
		log.Print("recieved signal to cancel/shutdown app")
		ctx2, _ := context.WithTimeout(shutdownContext, 10*time.Second)
		err := srv.Shutdown(ctx2)
		if err != nil {
			log.Printf("unable to shutdown? %+v", err)
		}
		shutdownComplete()
	}()

	log.Printf("listening on %q", port)
	err = srv.ListenAndServe()
	<-shutdownContext.Done()
	log.Printf("exiting: %+v", err)
}

// Should this take context instead of returning a channel?
func runBackground(ctx context.Context, f func(), duration time.Duration) {
	ticker := time.NewTicker(duration)

	log.Print("starting background updater...")

	go func() {
		f()
		for {
			select {
			case <-ticker.C:
				f()

			case <-ctx.Done():
				log.Print("recieved signal to cancel/shutdown background")
				ticker.Stop()
				return
			}
		}
	}()
}
