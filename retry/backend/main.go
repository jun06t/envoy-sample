package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	h := &handler{}
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/", genInstrumentChain("alive", h.alive))
	http.Handle("/hello", genInstrumentChain("hello", h.hello))
	http.ListenAndServe(":8000", nil)
}

type handler struct {
}

func (h *handler) alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive")
}

func (h *handler) hello(w http.ResponseWriter, r *http.Request) {
	dur := rand.Intn(1000)
	time.Sleep(time.Duration(dur) * time.Millisecond) // 処理を表現するためのsleep
	n := rand.Intn(4)                                 // エラーレスポンスを返すためのランダム値
	switch n {
	case 0:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello World")
	case 1:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Bad Request")
	case 2:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	case 3:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Server Unavailable")
	}
}
