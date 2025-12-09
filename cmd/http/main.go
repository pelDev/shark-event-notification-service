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

func main() {
	cfg := config.LoadConfig()

	// Renderer
	renderer, err := templates.NewGoTemplateRenderer(templates.Files)
	if err != nil {
		log.Fatalf("template init error: %v", err)
	}

	auth := smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.SMTPHost)

	emailProvider := providers.NewEmailProvider(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.From,
		renderer,
		auth,
	)

	// HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/send/email", func(w http.ResponseWriter, r *http.Request) {
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
	})

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
