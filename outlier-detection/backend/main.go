package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var errorRate atomic.Int32

func main() {
	instanceName := os.Getenv("INSTANCE_NAME")
	if instanceName == "" {
		instanceName = "unknown"
	}

	initialErrorRate := 0
	if tmp, err := strconv.Atoi(os.Getenv("ERROR_RATE")); err == nil {
		initialErrorRate = tmp
	}
	errorRate.Store(int32(initialErrorRate))

	h := &handler{
		instanceName: instanceName,
	}

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", genInstrumentChain("root", h.root))
	http.Handle("/health", genInstrumentChain("health", h.health))
	http.Handle("/error-rate", genInstrumentChain("error-rate", h.setErrorRate))

	log.Printf("Starting server %s on :8080 with error rate %d%%\n", instanceName, initialErrorRate)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type handler struct {
	instanceName string
}

func (h *handler) root(w http.ResponseWriter, r *http.Request) {
	rate := errorRate.Load()
	if rate > 0 && rand.IntN(100) < int(rate) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error from %s (error rate: %d%%)\n", h.instanceName, rate)
		return
	}
	fmt.Fprintf(w, "Hello from %s\n", h.instanceName)
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK from %s\n", h.instanceName)
}

func (h *handler) setErrorRate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		rateStr := r.URL.Query().Get("rate")
		rate, err := strconv.Atoi(rateStr)
		if err != nil || rate < 0 || rate > 100 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid error rate: %s (must be 0-100)\n", rateStr)
			return
		}
		errorRate.Store(int32(rate))
		log.Printf("Error rate changed to %d%%\n", rate)
		fmt.Fprintf(w, "Error rate set to %d%% for %s\n", rate, h.instanceName)
		return
	}

	fmt.Fprintf(w, "Current error rate: %d%% for %s\n", errorRate.Load(), h.instanceName)
}
