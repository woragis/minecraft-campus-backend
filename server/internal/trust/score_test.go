package trust

import "testing"

func TestApplyDeltaClamp(t *testing.T) {
	if got := ApplyDelta(100, 10); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
	if got := ApplyDelta(10, -20); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestDeltaForEvent(t *testing.T) {
	if got := DeltaForEvent("confirmed_report"); got != -15 {
		t.Fatalf("expected -15, got %d", got)
	}
}

func TestComputeSponsorScore(t *testing.T) {
	score := ComputeSponsorScore([]int{90, 80}, 0)
	if score != 85 {
		t.Fatalf("expected 85, got %d", score)
	}
	score = ComputeSponsorScore([]int{90, 80}, 1)
	if score != 75 {
		t.Fatalf("expected 75, got %d", score)
	}
	score = ComputeSponsorScore(nil, 0)
	if score != 100 {
		t.Fatalf("expected 100, got %d", score)
	}
}
