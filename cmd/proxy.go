package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jaredreisinger/asp"
	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "run a proxy server for specific device IDs",
	Run:   proxy,
}

// proxyConfig is the configuration for the HTTP proxy that sensorpush-proxy
// provides.
type proxyConfig struct {
	Port         string            `asp:"port,p,,port for the proxy to listen on"`
	Sensors      map[string]string `asp:"sensors,s,,sensors to proxy; use 'symbolicName=id,otherName=another-id'\nstyle formatting for the values"`
	UpdatePeriod time.Duration     `asp:"update-period,u,,duration between updates"`
}

type proxyCmdConfig struct {
	RootConfig `mapstructure:",squash"`
	Proxy      proxyConfig
}

var proxyDefaults = proxyCmdConfig{
	Proxy: proxyConfig{
		Port:         ":5375",
		UpdatePeriod: 5 * time.Minute,
	},
}

func init() {
	err := asp.Attach(proxyCmd, proxyDefaults, aspOptions...)
	cobra.CheckErr(err)

	rootCmd.AddCommand(proxyCmd)
}

func proxy(cmd *cobra.Command, args []string) {
	cfg, err := asp.Get[proxyCmdConfig](cmd)
	if err != nil {
		log.Fatalf("unable to get config: %+v", err)
	}
	// log.Printf("got config: %#v", cfg)

	user := cfg.SensorPush.Username
	pass := cfg.SensorPush.Password
	port := cfg.Proxy.Port
	sensors := cfg.Proxy.Sensors
	updatePeriod := cfg.Proxy.UpdatePeriod

	// ensure port has a ":" prefix?...
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	client, err := sensorpush.NewClient(user, pass)
	if err != nil {
		log.Fatalf("unable to create client: %+v", err)
	}

	// TODO: should lastSamples/SuccessfulCall really be channels?  We don't
	// really want to serve/read from them *while* updates are happening...
	lastSamples := make(map[string]sensorpush.Sample, len(sensors))
	lastSuccessfulCall := time.Now()

	// Maps are not safe for r/w concurrency:
	// https://go.dev/blog/maps#concurrency
	mutex := &sync.RWMutex{}

	appCtx, appCancel := context.WithCancel(context.Background())

	updater := func() {
		// ensure updates happen with a write-lock!
		mutex.Lock()
		defer mutex.Unlock()

		for key, id := range sensors {
			// log.Printf("getting last sample for %q...", key)
			sample, err := client.LastSample(id)
			if err != nil {
				log.Printf("unable to get sample for %q: %+v", key, err)

				// if we're 5 times past the last successful call, it's time to
				// cancel and exit, and let a new process/container start up...
				// but to be fair, we don't really have an expectation that
				// doing so will fix things, do we?
				if time.Since(lastSuccessfulCall) > (5 * updatePeriod) {
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

	runBackground(appCtx, updater, updatePeriod)

	// Ensures we read/marshal the map safely, protected against an incoming
	// update.
	getSamplesJson := func() ([]byte, error) {
		mutex.RLock()
		defer mutex.RUnlock()

		return json.Marshal(lastSamples)
	}
	// Now spin up a web server to serve the sample data...
	http.HandleFunc("/sensors", func(w http.ResponseWriter, req *http.Request) {
		// If there's an Origin header, send back Access-Control-Allow-Origin
		origin := req.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		b, err := getSamplesJson()
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
		log.Print("received signal to cancel/shutdown app")
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
				log.Print("received signal to cancel/shutdown background")
				ticker.Stop()
				return
			}
		}
	}()
}
