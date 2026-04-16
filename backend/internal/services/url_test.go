package services

import (
	"errors"
	"testing"
	"url-shortener/internal/storage"
)

func TestShortenNormalizesAndDeduplicates(t *testing.T) {
	service := New(storage.New())

	first, err := service.Shorten("example.com/path")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if first.LongURL != "https://example.com/path" {
		t.Fatalf("expected normalized https URL, got %q", first.LongURL)
	}

	second, err := service.Shorten("https://example.com/path")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if first.ShortCode != second.ShortCode {
		t.Fatalf("expected identical short code for duplicate URL, got %q and %q", first.ShortCode, second.ShortCode)
	}
}

func TestShortenRejectsInvalidURL(t *testing.T) {
	service := New(storage.New())

	_, err := service.Shorten("not a url @")
	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestShortenHandlesCodeCollisions(t *testing.T) {
	store := storage.New()
	sequence := []string{"aaaaaa", "aaaaaa", "bbbbbb"}
	index := 0

	service := newWithGenerator(store, func(_ int) string {
		if index >= len(sequence) {
			return "zzzzzz"
		}
		value := sequence[index]
		index++
		return value
	})

	first, err := service.Shorten("https://example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if first.ShortCode != "aaaaaa" {
		t.Fatalf("expected first code aaaaaa, got %q", first.ShortCode)
	}

	second, err := service.Shorten("https://example.org")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if second.ShortCode != "bbbbbb" {
		t.Fatalf("expected collision retry to pick bbbbbb, got %q", second.ShortCode)
	}
}

func TestResolveIncrementsClicks(t *testing.T) {
	service := New(storage.New())
	created, err := service.Shorten("https://clicks.test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = service.Resolve(created.ShortCode)
	if err != nil {
		t.Fatalf("expected resolve to succeed, got %v", err)
	}

	fetched, err := service.Get(created.ShortCode)
	if err != nil {
		t.Fatalf("expected get to succeed, got %v", err)
	}

	if fetched.Clicks != 1 {
		t.Fatalf("expected clicks to be 1, got %d", fetched.Clicks)
	}
}
