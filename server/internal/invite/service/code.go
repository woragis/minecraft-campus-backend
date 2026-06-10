package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const inviteCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const inviteCodeLength = 6

func generateInviteCode() (string, error) {
	b := make([]byte, inviteCodeLength)
	max := big.NewInt(int64(len(inviteCodeAlphabet)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("invite code random: %w", err)
		}
		b[i] = inviteCodeAlphabet[n.Int64()]
	}
	return "CW-" + string(b), nil
}
