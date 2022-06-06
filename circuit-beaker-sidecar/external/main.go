package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	loop := 10
	if tmp, _ := strconv.Atoi(os.Getenv("LOOP_COUNT")); tmp != 0 {
		loop = tmp
	}
	num := 30
	if tmp, _ := strconv.Atoi(os.Getenv("FIB_NUM")); tmp != 0 {
		num = tmp
	}

	h := &handler{
		time: 5 * time.Second,
		loop: loop,
		num:  num,
	}
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/", genInstrumentChain("alive", h.alive))
	http.Handle("/sleep", genInstrumentChain("sleep", h.sleep))
	http.Handle("/cpu", genInstrumentChain("cpu", h.cpu))
	http.ListenAndServe(":8002", nil)
}

type handler struct {
	time time.Duration
	loop int
	num  int
}

func (h *handler) alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive")
}

func (h *handler) sleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(h.time)
	fmt.Fprintf(w, "Sleep")
}

func (h *handler) cpu(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < h.loop; i++ {
		fib(h.num)
	}
	fmt.Fprintf(w, "CPU Bound")
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
