package httphandler

import (
	"net/http"
	"strconv"

	applicationdto "github.com/commitshark/notification-svc/internal/application/dto"
	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/domain/ports"
)

type NotificationHandler struct {
	notificationRepo ports.NotificationRepository
}

func NewNotificationHandler(notificationRepo ports.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
	}
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	req, err := parseListNotificationsRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request parameters", err)
		return
	}

	// Build filter from query params
	var status *domain.NotificationStatus
	if req.Status != "" {
		s := domain.NotificationStatus(req.Status)
		status = &s
	}

	var notifType *domain.NotificationType
	if req.Type != "" {
		t := domain.NotificationType(req.Type)
		notifType = &t
	}

	filter := domain.NotificationFilter{
		Status:      status,
		Type:        notifType,
		IsMarketing: req.IsMarketing,
		Query:       req.Query,
	}

	// Get paginated notifications from repository
	notifications, totalItems, err := h.notificationRepo.PaginatedList(ctx, req.Page, req.PageSize, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch notifications", err)
		return
	}

	// Create paginated response
	response := applicationdto.NewPaginatedResponse(applicationdto.ToNotificationDtos(notifications), req.Page, req.PageSize, int64(totalItems))

	// Return JSON response
	writeJSON(w, http.StatusOK, response)
}

func parseListNotificationsRequest(r *http.Request) (*applicationdto.ListNotificationsRequest, error) {
	req := &applicationdto.ListNotificationsRequest{
		Page:     1,
		PageSize: 20,
	}

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, err
		}
		if page > 0 {
			req.Page = page
		}
	}

	// Parse page size
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, err
		}
		if pageSize > 0 && pageSize <= 100 { // limit max page size to 100
			req.PageSize = pageSize
		}
	}

	// Parse optional filters
	req.Query = r.URL.Query().Get("q")
	req.Status = r.URL.Query().Get("status")
	req.Type = r.URL.Query().Get("type")

	isMarketingString := r.URL.Query().Get("is_marketing")

	var isMarketing *bool

	if isMarketingString == "" {
		isMarketing = nil
	} else if isMarketingString == "true" {
		v := true
		isMarketing = &v
	} else {
		v := false
		isMarketing = &v
	}

	req.IsMarketing = isMarketing

	return req, nil
}
