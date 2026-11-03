package handler

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type TimeResponse struct {
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
	Format   string `json:"format"`
}

type Handler struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/time", h.handleGetTime)

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../public/about.html")
	})

	mux.HandleFunc("GET /", h.handleStatic)
}

func (h *Handler) handleGetTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	now := time.Now()

	layout := "15:04:05"
	format := "24h"

	if r.URL.Query().Get("format") == "12" {
		layout = "03:04:05 pm"
		format = "12h"
	}

	data := TimeResponse{
		Time:     now.Format(layout),
		Timezone: now.Location().String(),
		Format:   format,
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("write time response", "err", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleStatic(w http.ResponseWriter, r *http.Request) {
	public := os.DirFS("../public")
	fileServer := http.FileServer(http.FS(public))

	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index.html"
	}

	info, err := fs.Stat(public, name)

	status := 0
	message := ""

	switch {
	case errors.Is(err, fs.ErrNotExist):
		status = http.StatusNotFound
		message = "Not found"
	case err != nil:
		h.log.Error("stat public file", "path", name, "err", err)
		status = http.StatusInternalServerError
		message = "Internal server error"
	case info.IsDir():
		status = http.StatusNotFound
		message = "Not found"
	}

	if status != 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": message,
		}); err != nil {
			h.log.Error("write error response", "err", err)
		}
		return
	}

	fileServer.ServeHTTP(w, r)
}
