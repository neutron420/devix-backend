package sanitize

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	// policy is the HTML sanitization policy for user-generated content.
	policy *bluemonday.Policy

	// strictPolicy strips ALL HTML tags.
	strictPolicy *bluemonday.Policy

	// filenameRegex matches invalid filename characters.
	filenameRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

	// multiSpaceRegex collapses multiple spaces.
	multiSpaceRegex = regexp.MustCompile(`\s+`)
)

func init() {
	// UGC policy — allows basic formatting but strips scripts, iframes, etc.
	policy = bluemonday.UGCPolicy()

	// Strict policy — strips everything
	strictPolicy = bluemonday.StrictPolicy()
}

// HTML sanitizes user-generated HTML content, allowing safe formatting tags.
func HTML(input string) string {
	return policy.Sanitize(input)
}

// StripHTML removes ALL HTML tags from input.
func StripHTML(input string) string {
	return strictPolicy.Sanitize(input)
}

// Text cleans a plain text input — trims whitespace and collapses multiple spaces.
func Text(input string) string {
	input = strings.TrimSpace(input)
	input = multiSpaceRegex.ReplaceAllString(input, " ")
	return input
}

// Filename sanitizes a filename for safe storage.
func Filename(input string) string {
	// Get just the filename (no path traversal)
	name := filepath.Base(input)

	// Get extension
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	// Replace invalid characters
	base = filenameRegex.ReplaceAllString(base, "_")

	// Ensure non-empty base
	if base == "" || base == "_" {
		base = "file"
	}

	// Lowercase extension
	ext = strings.ToLower(ext)

	return base + ext
}

// Slug sanitizes a string for use as a URL slug.
func Slug(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	input = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`[\s-]+`).ReplaceAllString(input, "-")
	input = strings.Trim(input, "-")
	return input
}
