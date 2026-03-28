package infrahttp

import (
	"net/http"
	"time"

	"github.com/commitshark/notification-svc/internal/domain/ports"
	httphandler "github.com/commitshark/notification-svc/internal/interfaces/http/handler"
	"github.com/commitshark/notification-svc/internal/interfaces/http/middlewares"
	"github.com/go-chi/chi"
	chi_middleware "github.com/go-chi/chi/middleware"
)

func NewRouter(
	notificationRepo ports.NotificationRepository,
) http.Handler {
	r := chi.NewRouter()

	// -------------------
	// Global middleware
	// -------------------
	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	r.Use(chi_middleware.Timeout(30 * time.Second))

	// -------------------
	// Handlers
	// -------------------
	handler := httphandler.NewNotificationHandler(notificationRepo)

	// -------------------
	// Middleware
	// -------------------
	authn := middlewares.NewAuthnMiddleware()

	// -------------------
	// Routes
	// -------------------

	r.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authn.RequireSession)
			r.Use(authn.RequireAdmin)

			r.Get("/", handler.ListNotifications)
		})
	})

	return r
}
