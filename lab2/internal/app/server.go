package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"lab2/internal/config"
	"lab2/internal/middleware"
	"lab2/internal/model"
	"lab2/internal/repository"
	"lab2/internal/service"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	messageRepository := repository.NewInMemoryMessageRepository([]model.Message{
		{ID: 1, Text: "hello"},
		{ID: 2, Text: "world"},
	})
	messageService := service.NewMessageService(messageRepository)
	server := &http.Server{Addr: cfg.HTTPAddress()}
	server.Handler = routes(server, messageService, cfg.ShutdownToken, cfg.ShutdownTimeout)
	server.ReadHeaderTimeout = cfg.ReadHeaderTimeout
	server.ReadTimeout = cfg.ReadTimeout
	server.WriteTimeout = cfg.WriteTimeout
	server.IdleTimeout = cfg.IdleTimeout

	log.Printf("server is running on %s", cfg.HTTPAddress())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func routes(server *http.Server, messageService *service.MessageService, shutdownToken string, shutdownTimeout time.Duration) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /messages", listMessagesHandler(messageService))
	mux.HandleFunc("POST /messages", createMessageHandler(messageService))
	mux.Handle("POST /admin/shutdown", middleware.RequireBearerToken(shutdownToken, http.HandlerFunc(shutdownHandler(server, shutdownTimeout))))

	return middleware.RecoverPanic(middleware.RequestLogging(mux))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func shutdownHandler(server *http.Server, shutdownTimeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"message": "server is shutting down"})

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Printf("shutdown error: %v", err)
			}
		}()
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		http.Error(w, "encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(body.Bytes())
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
