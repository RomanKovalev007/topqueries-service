package domain

import (
	"time"

	"github.com/google/uuid"
)

type TopEntry struct {
    Query string `json:"query"`
    Count int64  `json:"count"`
}

type SearchEvent struct {
	QueryText   string    `json:"query"`
	UserID      uuid.UUID `json:"user_id"`
	TimeRequest time.Time `json:"timestamp"`
}