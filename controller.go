package cantainer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/netip"
	"time"
)

type HTTP struct {
	server *http.Server
	core   *core
}

func NewHTTP(core *core) *HTTP {

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":20043",
		Handler: mux,
	}

	controller := HTTP{
		server: &server,
		core:   core,
	}

	mux.HandleFunc("GET /ip", controller.NetworkIP)

	return &controller
}

func (h *HTTP) Run(ctx context.Context) error {

	go func() {

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := h.server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}()

	return h.server.ListenAndServe()
}

type RegisterContainerHTTPRequest struct {
	Network netip.Prefix `json:"network"`
}

type RegisterContainerHTTPResponse struct {
	Address netip.Addr `json:"address"`
}

func (h *HTTP) NetworkIP(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)

	var request RegisterContainerHTTPRequest
	if err := dec.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addr, err := h.core.RegisterContainer(context.Background(), &request.Network)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(RegisterContainerHTTPResponse{Address: addr}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
