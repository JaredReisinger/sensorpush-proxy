package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(proxyCmd)
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "run a proxy server for specific device IDs",
	Run:   proxy,
}

func proxy(cmd *cobra.Command, args []string) {
	config := getConfig(cmd)
	// log.Printf("got config: %#v", config)

	user := config.SensorPush.Username
	pass := config.SensorPush.Password
	port := config.Proxy.Port
	deviceIDs := config.Proxy.DeviceIDs

	// TODO: check flags!

	// ensure port has a ":" prefix?...
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	client, err := sensorpush.NewClient(user, pass)
	if err != nil {
		log.Fatalf("unable to create client: %+v", err)
	}

	// TODO: should lastSamples/SuccessfulCall really be channels?  We don't
	// really want to server/read from them *while* updates are happening...
	lastSamples := make(map[string]sensorpush.Sample, len(deviceIDs))
	lastSuccessfulCall := time.Now()

	appCtx, appCancel := context.WithCancel(context.Background())

	updater := func() {
		for key, id := range deviceIDs {
			// log.Printf("getting last sample for %q...", key)
			sample, err := client.LastSample(id)
			if err != nil {
				log.Printf("unable to get sample for %q: %+v", key, err)
				// if we're X? past the last successful call, it's time to cancel
				// and exit, and let a new process/container start up
				if time.Since(lastSuccessfulCall) > (5 * time.Minute) {
					appCancel()
				}
			} else {
				log.Printf(
					"UPDATER: %s: %.1f°F, %.1f%%RH (%s)",
					sample.Observed.Local().Format(time.RFC3339),
					sample.Temperature,
					sample.Humidity,
					key,
				)
				lastSamples[key] = *sample
				lastSuccessfulCall = time.Now()
			}
		}
	}

	runBackground(appCtx, updater, time.Minute)

	// Now spin up a web server to serve the sample data...
	http.HandleFunc("/sensors", func(w http.ResponseWriter, req *http.Request) {
		// If there's an Origin header, send back Access-Control-Allow-Origin
		origin := req.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		b, err := json.Marshal(lastSamples)
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
		ctx2, cancel := context.WithTimeout(shutdownContext, 10*time.Second)
		defer cancel()
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
