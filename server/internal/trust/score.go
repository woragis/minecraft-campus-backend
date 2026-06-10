package trust

// ApplyDelta returns a trust score clamped to [0, 100].
func ApplyDelta(current, delta int) int {
	score := current + delta
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

// DeltaForEvent returns the trust delta for a known event type.
func DeltaForEvent(eventType string) int {
	switch eventType {
	case "confirmed_report":
		return -15
	case "rollback_applied":
		return -5
	case "probation_day_clean":
		return 1
	case "ban":
		return -50
	case "unban":
		return 0
	default:
		return 0
	}
}

// ComputeSponsorScore averages direct invitee trust and applies ban penalty.
func ComputeSponsorScore(inviteeTrustScores []int, recentInviteeBans int) int {
	if len(inviteeTrustScores) == 0 {
		return 100
	}
	sum := 0
	for _, s := range inviteeTrustScores {
		sum += s
	}
	avg := sum / len(inviteeTrustScores)
	score := avg - (recentInviteeBans * 10)
	return ApplyDelta(score, 0)
}
