package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	Version    string = "set-at-build-time"
	listenAddr string = ":8080"
	isReady    bool
	hostname   string
)

func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("%s %s home\n", Version, hostname))
}

func readiness(w http.ResponseWriter, r *http.Request) {
	if !isReady {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("%s %s is unready\n", Version, hostname))
		log.Printf("%s %s: readiness queried - replied unready", Version, hostname)
		return
	}

	io.WriteString(w, fmt.Sprintf("%s %s is ready\n", Version, hostname))
	log.Printf("%s %s: readiness queried - replied ready", Version, hostname)
}

func unready(w http.ResponseWriter, r *http.Request) {
	isReady = false
	io.WriteString(w, fmt.Sprintf("%s %s made unready\n", Version, hostname))
	log.Printf("%s %s: was made unready", Version, hostname)
}

func liveness(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("%s %s is alive\n", Version, hostname))
	log.Printf("%s %s: liveness queried - replied alive", Version, hostname)
}

func dead(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("%s %s about to crash\n", Version, hostname))
	panic(fmt.Sprintf("%s %s: i have crashed", Version, hostname))
}

func runServer(ctx context.Context, wg *sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/readiness", readiness)
	mux.HandleFunc("/unready", unready)
	mux.HandleFunc("/liveness", liveness)
	mux.HandleFunc("/dead", dead)

	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		log.Printf("%s %s: listening on %s\n", Version, hostname, listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("%s %s: Error: %v\n", Version, hostname, err)
		}
		wg.Done()
	}()

	<-ctx.Done()
	// Immediately start routing traffic elsewhere, and give time
	// for routing to adjust
	isReady = false
	time.Sleep(time.Second * 10)

	// Give existing connections 5 seconds to finish, then force close
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for existing connections to finish
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("%s %s: could not elegantly shutdown: %v", Version, hostname, err)
	}

	log.Printf("%s %s: server gracefully shutdown", Version, hostname)
	wg.Done()
}

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-c
	log.Printf("%s %s: Received signal: %v\n", Version, hostname, sig)

	// Trigger the cancellation of the context to initiate shutdown
	cancel()
}

func main() {
	hostname, _ = os.Hostname()
	isReady = true

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go runServer(ctx, &wg)
	handleSignals(cancel)

	wg.Wait()
}
