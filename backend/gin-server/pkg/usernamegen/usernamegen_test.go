package usernamegen

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateFormat(t *testing.T) {
	username := Generate()
	parts := strings.Split(username, "-")

	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d in username '%s'", len(parts), username)
	}

	// Check last part is a valid timestamp (Unix seconds)
	if _, err := time.Parse("2006-01-02T15:04:05Z07:00", parts[3]); err == nil {
		t.Errorf("expected Unix timestamp, got RFC3339-like timestamp: %s", parts[3])
	}
}

func TestGenerateUniqueness(t *testing.T) {
	username1 := Generate()
	time.Sleep(time.Millisecond * 10) // Ensure timestamp changes
	username2 := Generate()

	if username1 == username2 {
		t.Error("expected unique usernames, but got duplicates")
	}
}

func TestGenerateValidWords(t *testing.T) {
	username := Generate()
	parts := strings.Split(username, "-")

	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d", len(parts))
	}

	adverb := parts[0]
	adjective := parts[1]
	animal := parts[2]

	if !contains(adverbs, adverb) {
		t.Errorf("adverb '%s' not found in list", adverb)
	}
	if !contains(adjectives, adjective) {
		t.Errorf("adjective '%s' not found in list", adjective)
	}
	if !contains(animals, animal) {
		t.Errorf("animal '%s' not found in list", animal)
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
