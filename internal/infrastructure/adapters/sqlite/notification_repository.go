package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/commitshark/notification-svc/internal/domain"
	"github.com/commitshark/notification-svc/internal/utils"
	_ "modernc.org/sqlite"
)

type SQLiteNotificationRepository struct {
	db *sql.DB
}

func NewSQLiteNotificationRepository(dbPath string) (*SQLiteNotificationRepository, error) {
	// SQLite with WAL mode for better concurrency
	dsn := fmt.Sprintf("%s?_journal=WAL&_timeout=5000&_fk=true", dbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err := optimizeSQLite(db); err != nil {
		return nil, fmt.Errorf("failed to optimize db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &SQLiteNotificationRepository{db: db}, nil
}

// Add these to SQLite setup
func optimizeSQLite(db *sql.DB) error {
	_, err := db.Exec(`
        PRAGMA journal_mode = WAL;
        PRAGMA synchronous = NORMAL;
        PRAGMA cache_size = -10000;  -- 10MB cache
        PRAGMA mmap_size = 30000000000;
        PRAGMA busy_timeout = 5000;
        PRAGMA auto_vacuum = INCREMENTAL;
    `)
	return err
}

func createTables(db *sql.DB) error {
	// Notifications table
	notificationTable := `
	CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		recipient_id TEXT NOT NULL,
		recipient_email TEXT,
		recipient_phone TEXT,
		recipient_device TEXT,
		title TEXT NOT NULL,
		body TEXT,
		data TEXT, -- JSON data
		html TEXT,
		template TEXT,
		status TEXT NOT NULL,
		provider_response TEXT,
		created_at DATETIME NOT NULL,
		sent_at DATETIME,
		retry_count INTEGER DEFAULT 0,
		max_retries INTEGER DEFAULT 3,
		version INTEGER DEFAULT 1,
		CHECK (type IN ('EMAIL', 'SMS', 'PUSH')),
		CHECK (status IN ('PENDING', 'SENT', 'FAILED', 'DELIVERED'))
	);
	`

	// Indexes for query performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status, created_at)",
		"CREATE INDEX IF NOT EXISTS idx_notifications_recipient ON notifications(recipient_id)",
		"CREATE INDEX IF NOT EXISTS idx_notifications_retry ON notifications(status, retry_count, created_at) WHERE status = 'FAILED'",
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(notificationTable); err != nil {
		return fmt.Errorf("failed to create notifications table: %w", err)
	}

	for _, index := range indexes {
		if _, err := tx.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return tx.Commit()
}

func (r *SQLiteNotificationRepository) Save(ctx context.Context, notification *domain.Notification) error {
	dataJSON, err := json.Marshal(notification.Content.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := `
	INSERT INTO notifications (
		id, type, recipient_id, recipient_email, recipient_phone,
		recipient_device, title, body, data, status, provider_response,
		created_at, sent_at, retry_count, max_retries, html, template, version
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status,
		provider_response = excluded.provider_response,
		sent_at = excluded.sent_at,
		retry_count = excluded.retry_count,
		version = version + 1
	WHERE version = ?
	`

	var sentAt interface{}
	if notification.SentAt != nil {
		sentAt = notification.SentAt.Format(time.RFC3339)
	}

	args := []interface{}{
		notification.ID,
		string(notification.Type),
		notification.Recipient.ID,
		notification.Recipient.Email,
		notification.Recipient.Phone,
		notification.Recipient.DeviceID,
		notification.Content.Title,
		notification.Content.Body,
		string(dataJSON),
		string(notification.Status),
		notification.ProviderResponse,
		notification.CreatedAt.Format(time.RFC3339),
		sentAt,
		notification.RetryCount,
		notification.MaxRetries,
		notification.Version,
		notification.Content.HTML,
		notification.Content.Template,
		notification.Version - 1, // For optimistic locking
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("optimistic locking conflict for notification %s", notification.ID)
	}

	return nil
}

func (r *SQLiteNotificationRepository) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	query := `
	SELECT 
		id, type, recipient_id, recipient_email, recipient_phone,
		recipient_device, html, template,
		title, body, data, status, provider_response,
		created_at, sent_at, retry_count, max_retries, version
	FROM notifications 
	WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var n domain.Notification
	var recipientID string
	var recipientEmail, recipientPhone, recipientDevice, html, template sql.NullString
	var title, statusStr, typeStr string
	var body, dataJSON sql.NullString
	var providerResponse sql.NullString
	var createdAtStr string
	var sentAtStr sql.NullString
	var data map[string]interface{}

	err := row.Scan(
		&n.ID, &typeStr, &recipientID, &recipientEmail, &recipientPhone, &recipientDevice,
		&html, &template,
		&title, &body, &dataJSON, &statusStr, &providerResponse,
		&createdAtStr, &sentAtStr, &n.RetryCount, &n.MaxRetries, &n.Version,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notification not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan notification: %w", err)
	}

	// Parse JSON data
	if dataJSON.Valid && dataJSON.String != "" {
		if err := json.Unmarshal([]byte(dataJSON.String), &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	var sentAt *time.Time
	if sentAtStr.Valid {
		t, err := time.Parse(time.RFC3339, sentAtStr.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sent_at: %w", err)
		}
		sentAt = &t
	}

	// Build domain objects
	recipient := domain.Recipient{
		ID:       recipientID,
		Email:    utils.SqlNullableString(recipientEmail),
		Phone:    utils.SqlNullableString(recipientPhone),
		DeviceID: utils.SqlNullableString(recipientDevice),
	}

	content := domain.Content{
		Title:    title,
		Body:     utils.SqlNullableString(body),
		Data:     &data,
		HTML:     utils.SqlNullableString(html),
		Template: utils.SqlNullableString(template),
	}

	n.Type = domain.NotificationType(typeStr)
	n.Status = domain.NotificationStatus(statusStr)
	n.Recipient = recipient
	n.Content = content
	n.ProviderResponse = providerResponse.String
	n.CreatedAt = createdAt
	n.SentAt = sentAt

	return &n, nil
}

func (r *SQLiteNotificationRepository) FindPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `
	SELECT id FROM notifications 
	WHERE status IN ('PENDING', 'FAILED') 
		AND retry_count < max_retries
	ORDER BY created_at ASC
	LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending notifications: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	// Load each notification (could be optimized with a JOIN)
	var notifications []*domain.Notification
	for _, id := range ids {
		notification, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *SQLiteNotificationRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.NotificationStatus,
	providerResponse string,
) error {
	query := `
	UPDATE notifications 
	SET status = ?, 
		provider_response = ?,
		sent_at = CASE WHEN ? = 'SENT' THEN CURRENT_TIMESTAMP ELSE sent_at END,
		version = version + 1
	WHERE id = ? AND version = ?
	`

	// Get current version first
	var currentVersion int
	versionQuery := `SELECT version FROM notifications WHERE id = ?`
	err := r.db.QueryRowContext(ctx, versionQuery, id).Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get notification version: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		string(status),
		providerResponse,
		string(status),
		id,
		currentVersion,
	)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("optimistic locking conflict for notification %s", id)
	}

	return nil
}

func (r *SQLiteNotificationRepository) IncrementRetryCount(ctx context.Context, id string) error {
	query := `
	UPDATE notifications 
	SET retry_count = retry_count + 1,
		version = version + 1
	WHERE id = ? AND version = ?
	`

	var currentVersion int
	versionQuery := `SELECT version FROM notifications WHERE id = ?`
	err := r.db.QueryRowContext(ctx, versionQuery, id).Scan(&currentVersion)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, id, currentVersion)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("optimistic locking conflict")
	}

	return nil
}

func (r *SQLiteNotificationRepository) Close() error {
	return r.db.Close()
}
