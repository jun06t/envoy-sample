package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	h := &handler{
		url: os.Getenv("ENDPOINT"),
	}
	http.HandleFunc("/", h.alive)
	http.HandleFunc("/test", h.external)
	http.HandleFunc("/httpbin", h.httpbin)
	http.ListenAndServe(":8001", nil)
}

type handler struct {
	url string
}

func (h *handler) alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive")
}

func (h *handler) external(w http.ResponseWriter, r *http.Request) {
	cli := http.DefaultClient
	resp, err := cli.Get(h.url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 503:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "external service returns 503")
	case 504:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "external service returns 504")
	default:
		fmt.Fprintf(w, "Succeed in calling external api")
	}
}

func (h *handler) httpbin(w http.ResponseWriter, r *http.Request) {
	cli := http.DefaultClient
	resp, err := cli.Get("http://sidecar:8083/status/200")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 503:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "httpbin returns 503")
	case 504:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "httpbin returns 504")
	default:
		fmt.Fprintf(w, "Succeed in calling httpbin")
	}
}
