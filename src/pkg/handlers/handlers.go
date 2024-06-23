package handlers

import (
	"fmt"
	"net/http"
)

type Handlers struct {
}

func (h *Handlers) Healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func (h *Handlers) ServeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome! The configuration - Every: %d %s\n")
}
