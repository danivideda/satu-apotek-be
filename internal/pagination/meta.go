package pagination

import "math"

const (
	maxSize     = 100
	defaultSize = 50
)

type Meta struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// NewMeta builds pagination meta from 1-based page, size, and total items.
// size is the limit (page size). Invalid inputs are normalized.
func NewMeta(page, size int, total int64) Meta {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = defaultSize 
	}
	if size > maxSize {
		size = maxSize
	}

	var totalPages int
	if total > 0 && size > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(size)))
	}

	return Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
