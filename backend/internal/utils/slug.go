package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

// GenerateSlugFromName generates a unique slug in the format: {waitlist-name}-{6 char alphanumeric}
func GenerateSlugFromName(name string) string {
	// Convert name to lowercase and replace spaces/special chars with hyphens
	slug := strings.ToLower(name)
	// Remove any characters that aren't alphanumeric, spaces, or hyphens
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Replace multiple consecutive hyphens with single hyphen
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Generate 6-character alphanumeric suffix
	suffix := generateRandomString(6)

	return slug + "-" + suffix
}

// generateRandomString generates a random alphanumeric string of given length
func generateRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
