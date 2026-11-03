package main

import (
	"go_server/internal/handler"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	log := slog.Default()
	h := handler.New(log)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	srv := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err := srv.ListenAndServe()

	if err != nil {
		slog.Error("server stopped with error", "err", err)
	}
}
