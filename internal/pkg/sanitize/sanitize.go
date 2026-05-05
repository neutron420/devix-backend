package sanitize

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var (
	policy *bluemonday.Policy

	strictPolicy *bluemonday.Policy

	filenameRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

	multiSpaceRegex = regexp.MustCompile(`\s+`)
)

func init() {

	policy = bluemonday.UGCPolicy()

	strictPolicy = bluemonday.StrictPolicy()
}

func HTML(input string) string {
	return policy.Sanitize(input)
}

func StripHTML(input string) string {
	return strictPolicy.Sanitize(input)
}

func Text(input string) string {
	input = strings.TrimSpace(input)
	input = multiSpaceRegex.ReplaceAllString(input, " ")
	return input
}

func Filename(input string) string {

	name := filepath.Base(input)

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	base = filenameRegex.ReplaceAllString(base, "_")

	if base == "" || base == "_" {
		base = "file"
	}

	ext = strings.ToLower(ext)

	return base + ext
}

func Slug(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	input = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`[\s-]+`).ReplaceAllString(input, "-")
	input = strings.Trim(input, "-")
	return input
}
