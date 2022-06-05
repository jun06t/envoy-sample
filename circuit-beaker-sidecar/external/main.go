package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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
		loop: loop,
		num:  num,
	}
	http.HandleFunc("/", h.alive)
	http.HandleFunc("/test", h.heavy)
	http.ListenAndServe(":8002", nil)
}

type handler struct {
	loop int
	num  int
}

func (h *handler) alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive")
}

func (h *handler) heavy(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < h.loop; i++ {
		fib(h.num)
	}
	fmt.Fprintf(w, "Hello, World")
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
