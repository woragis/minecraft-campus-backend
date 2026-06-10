package territory

import (
	"time"
)

// MaxClaimArea returns the maximum horizontal claim area (blocks²) for a player.
func MaxClaimArea(accountCreated time.Time, trustScore int) int {
	days := time.Since(accountCreated).Hours() / 24
	base := 1000
	switch {
	case days >= 180:
		base = 10000
	case days >= 30:
		base = 3000
	case days >= 7:
		base = 1500
	}
	if trustScore < 50 {
		base = base / 2
	}
	if base < 100 {
		return 100
	}
	return base
}

// HorizontalArea computes (maxX-minX+1)*(maxZ-minZ+1).
func HorizontalArea(minX, maxX, minZ, maxZ int) int {
	w := maxX - minX + 1
	d := maxZ - minZ + 1
	if w <= 0 || d <= 0 {
		return 0
	}
	return w * d
}

// NormalizeBounds ensures min <= max for each axis.
func NormalizeBounds(minX, maxX, minZ, maxZ int) (int, int, int, int) {
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minZ > maxZ {
		minZ, maxZ = maxZ, minZ
	}
	return minX, maxX, minZ, maxZ
}

func ValidZone(zone string) bool {
	switch zone {
	case "urban", "rural", "industrial", "historic":
		return true
	default:
		return false
	}
}
