package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	config "github.com/commitshark/notification-svc/internal"
	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/providers"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/templates"
)

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvIntOrDefault(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var parsed int
		fmt.Sscanf(v, "%d", &parsed)
		return parsed
	}
	return fallback
}

func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expected := os.Getenv("API_KEY")
		if expected == "" {
			http.Error(w, "API_KEY not configured on server", http.StatusInternalServerError)
			return
		}

		received := r.Header.Get("X-API-Key")
		if received == "" || received != expected {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	cfg := config.LoadConfig()

	// Renderer
	renderer, err := templates.NewGoTemplateRenderer(templates.Files)
	if err != nil {
		log.Fatalf("template init error: %v", err)
	}

	auth := smtp.PlainAuth(
		"",
		getEnvOrDefault("SMTP_USERNAME", cfg.Email.Username),
		getEnvOrDefault("SMTP_PASSWORD", cfg.Email.Password),
		getEnvOrDefault("SMTP_HOST", cfg.Email.SMTPHost),
	)

	emailProvider := providers.NewEmailProvider(
		getEnvOrDefault("SMTP_HOST", cfg.Email.SMTPHost),
		getEnvIntOrDefault("SMTP_PORT", cfg.Email.SMTPPort),
		getEnvOrDefault("SMTP_USERNAME", cfg.Email.Username),
		getEnvOrDefault("SMTP_PASSWORD", cfg.Email.Password),
		getEnvOrDefault("SMTP_FROM", cfg.Email.From),
		renderer,
		auth,
	)

	// HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/send/email", requireAPIKey(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var n domain.Notification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if !emailProvider.Supports(n.Type) {
			http.Error(w, "Notification Type not supported: "+string(n.Type), http.StatusBadRequest)
			return
		}

		providerResponse, err := emailProvider.Send(&n)
		if err != nil {
			http.Error(w, "failed to send notification: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":   "sent",
			"id":       n.ID,
			"provider": providerResponse,
		})

		fmt.Println("Email sent →", providerResponse)
	}))

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("HTTP server running on : %s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
