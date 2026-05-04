package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// CursorParams holds parsed cursor pagination parameters.
type CursorParams struct {
	Cursor    string
	Limit     int
	Direction string // "next" or "prev"
}

// DecodedCursor holds the decoded cursor values.
type DecodedCursor struct {
	CreatedAt time.Time
	ID        string
}

// DefaultLimit is the default page size.
const DefaultLimit = 20

// MaxLimit is the maximum page size.
const MaxLimit = 100

// NormalizeLimit clamps the requested limit to valid bounds.
func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

// EncodeCursor creates a base64-encoded cursor string from a timestamp and ID.
func EncodeCursor(createdAt time.Time, id string) string {
	raw := fmt.Sprintf("%s|%s", createdAt.UTC().Format(time.RFC3339Nano), id)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor parses a base64-encoded cursor string.
func DecodeCursor(cursor string) (*DecodedCursor, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format")
	}

	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor data")
	}

	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid cursor timestamp")
	}

	return &DecodedCursor{
		CreatedAt: createdAt,
		ID:        parts[1],
	}, nil
}
