package pagination

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

const cursorPrefix = "v1:"

var cursorSecret []byte

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

func SetCursorSecret(secret string) {
	if secret == "" {
		cursorSecret = nil
		return
	}
	sum := sha256.Sum256([]byte(secret))
	cursorSecret = sum[:]
}

func EncodeCursor(createdAt time.Time, id string) string {
	raw := fmt.Sprintf("%s|%s", createdAt.UTC().Format(time.RFC3339Nano), id)
	if len(cursorSecret) == 0 {
		return base64.URLEncoding.EncodeToString([]byte(raw))
	}

	encoded, err := encryptCursor(raw)
	if err != nil {
		return base64.URLEncoding.EncodeToString([]byte(raw))
	}

	return cursorPrefix + encoded
}

func DecodeCursor(cursor string) (*DecodedCursor, error) {
	if cursor == "" {
		return nil, nil
	}

	var raw string
	if strings.HasPrefix(cursor, cursorPrefix) {
		if len(cursorSecret) == 0 {
			return nil, fmt.Errorf("cursor secret not configured")
		}
		decoded, err := decryptCursor(strings.TrimPrefix(cursor, cursorPrefix))
		if err != nil {
			return nil, fmt.Errorf("invalid cursor format")
		}
		raw = decoded
	} else {
		decoded, err := base64.URLEncoding.DecodeString(cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor format")
		}
		raw = string(decoded)
	}

	parts := strings.SplitN(raw, "|", 2)
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

func encryptCursor(plain string) (string, error) {
	block, err := aes.NewCipher(cursorSecret)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plain), nil)
	payload := append(nonce, ciphertext...)
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func decryptCursor(encoded string) (string, error) {
	if len(cursorSecret) == 0 {
		return "", errors.New("cursor secret not configured")
	}
	block, err := aes.NewCipher(cursorSecret)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(payload) < gcm.NonceSize() {
		return "", errors.New("invalid cursor payload")
	}

	nonce := payload[:gcm.NonceSize()]
	ciphertext := payload[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
