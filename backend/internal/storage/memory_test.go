package storage

import (
	"testing"
	"time"
)

func TestGetAllReturnsNewestFirst(t *testing.T) {
	store := New()

	older := store.Set("first11", "https://first.example")
	time.Sleep(10 * time.Millisecond)
	newer := store.Set("second2", "https://second.example")

	all := store.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 URLs, got %d", len(all))
	}

	if all[0].ShortCode != newer.ShortCode {
		t.Fatalf("expected newest URL first, got %q", all[0].ShortCode)
	}

	if all[1].ShortCode != older.ShortCode {
		t.Fatalf("expected older URL second, got %q", all[1].ShortCode)
	}
}
