package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type CursorParams struct {
	Cursor    string
	Limit     int
	Direction string
}

type DecodedCursor struct {
	CreatedAt time.Time
	ID        string
}

const DefaultLimit = 20

const MaxLimit = 100

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

func EncodeCursor(createdAt time.Time, id string) string {
	raw := fmt.Sprintf("%s|%s", createdAt.UTC().Format(time.RFC3339Nano), id)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

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
