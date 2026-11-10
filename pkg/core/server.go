package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
)

func NewServer(
	address string,
	sensorPushUsername string,
	sensorPushPassword string,
	sensors map[string]string,
	updatePeriod time.Duration,
) (*Server, error) {
	// We used to unilaterally prefix the address (port) with ":" if it was
	// missing, but it's completely valid to pass something like "0.0.0.0:1234",
	// so looking for a prefix isn't really the right thing to do. That said, we
	// also don't want to break compatibility (at least until a major release)
	// and fail if just "1234" is sent, so we check to see if the address parses
	// as a number, and if so, prefix with ":".
	_, err := strconv.Atoi(address)
	if err == nil && !strings.HasPrefix(address, ":") {
		address = fmt.Sprintf(":%s", address)
	}

	client, err := sensorpush.NewClient(sensorPushUsername, sensorPushPassword)
	if err != nil {
		// log.Fatalf("unable to create client: %+v", err)
		return nil, err
	}

	appCtx, appCancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	s := &Server{
		server:             httpServer,
		client:             client,
		sensors:            maps.Clone(sensors),
		updatePeriod:       updatePeriod,
		lastSamples:        make(map[string]sensorpush.Sample, len(sensors)),
		mutex:              &sync.RWMutex{},
		lastSuccessfulCall: time.Now(),
		appCtx:             appCtx,
		appCancel:          appCancel,
	}

	mux.HandleFunc("/sensors", s.sensorsHandler)

	return s, nil
}

type Server struct {
	server *http.Server

	client       *sensorpush.Client
	sensors      map[string]string
	updatePeriod time.Duration

	lastSamples        map[string]sensorpush.Sample
	mutex              *sync.RWMutex // maps are not safe for r/w concurrency
	lastSuccessfulCall time.Time

	appCtx    context.Context
	appCancel context.CancelFunc
}

func (s *Server) ListenAndServe() {
	s.runBackgroundUpdates()
	shutdownContext := s.runShutdownHandler()

	log.Printf("listening on %q", s.server.Addr)
	err := s.server.ListenAndServe()
	<-shutdownContext.Done()
	log.Printf("exiting: %+v", err)
}

func (s *Server) sensorsHandler(w http.ResponseWriter, req *http.Request) {
	// If there's an Origin header, send back Access-Control-Allow-Origin
	origin := req.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := s.getSamplesJson()
	if err != nil {
		io.WriteString(w, "{ error: \"no data?\" }")
		return
	}
	w.Write(b)
}

// Ensures we read/marshal the map safely, protected against an incoming
// update.
func (s *Server) getSamplesJson() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return json.Marshal(s.lastSamples)
}

func (s *Server) updater() {
	// ensure updates happen with a write-lock!
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, id := range s.sensors {
		// log.Printf("getting last sample for %q...", key)
		sample, err := s.client.LastSample(id)
		if err != nil {
			log.Printf("unable to get sample for %q: %+v", key, err)

			// if we're 5 times past the last successful call, it's time to
			// cancel and exit, and let a new process/container start up...
			// but to be fair, we don't really have an expectation that
			// doing so will fix things, do we?
			if time.Since(s.lastSuccessfulCall) > (5 * s.updatePeriod) {
				s.appCancel()
			}
		} else {
			log.Printf(
				"UPDATER: %s: %.1f°F, %.1f%%RH (%s)",
				sample.Observed.Local().Format(time.RFC3339),
				sample.Temperature,
				sample.Humidity,
				key,
			)
			s.lastSamples[key] = *sample
			s.lastSuccessfulCall = time.Now()
		}
	}
}

// Should this take context instead of returning a channel?
func (s *Server) runBackgroundUpdates() {

	log.Print("starting background updater...")

	ticker := time.NewTicker(s.updatePeriod)
	go func() {
		s.updater()
		for {
			select {
			case <-ticker.C:
				s.updater()

			case <-s.appCtx.Done():
				log.Print("received signal to cancel/shutdown background")
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Server) runShutdownHandler() context.Context {
	// also run a one-off background goroutine that closes the server when the
	// context is canceled...
	shutdownContext, shutdownComplete := context.WithCancel(context.Background())
	go func() {
		<-s.appCtx.Done()
		log.Print("received signal to cancel/shutdown app")
		ctx, cancel := context.WithTimeout(shutdownContext, 10*time.Second)
		defer cancel()
		err := s.server.Shutdown(ctx)
		if err != nil {
			log.Printf("unable to shutdown? %+v", err)
		}
		shutdownComplete()
	}()

	return shutdownContext
}
