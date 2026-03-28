package applicationdto

// PaginatedResponse is a generic pagination response
type PaginatedResponse[T any] struct {
	Data        []T   `json:"items"`
	TotalItems  int64 `json:"total"`
	CurrentPage int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, currentPage, pageSize int, totalItems int64) PaginatedResponse[T] {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize != 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}

	return PaginatedResponse[T]{
		Data:        data,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
	}
}
