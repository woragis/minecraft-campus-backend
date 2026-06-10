package territory

import (
	"testing"
	"time"
)

func TestMaxClaimAreaNewAccount(t *testing.T) {
	created := time.Now().UTC()
	if got := MaxClaimArea(created, 100); got != 1000 {
		t.Fatalf("expected 1000, got %d", got)
	}
}

func TestMaxClaimAreaSixMonths(t *testing.T) {
	created := time.Now().UTC().Add(-200 * 24 * time.Hour)
	if got := MaxClaimArea(created, 100); got != 10000 {
		t.Fatalf("expected 10000, got %d", got)
	}
}

func TestHorizontalArea(t *testing.T) {
	if got := HorizontalArea(0, 9, 0, 9); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}
