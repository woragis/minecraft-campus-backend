package webauth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrLinkCodeInvalid = errors.New("link code invalid or expired")
	ErrSessionInvalid  = errors.New("session invalid or expired")
)

type LinkPayload struct {
	PlayerID      uuid.UUID
	MinecraftUUID uuid.UUID
	Username      string
}

type Session struct {
	Token         string
	PlayerID      uuid.UUID
	MinecraftUUID uuid.UUID
	Username      string
	ExpiresAt     time.Time
}

type Service struct {
	mu          sync.Mutex
	links       map[string]linkEntry
	sessions    map[string]sessionEntry
	linkTTL     time.Duration
	sessionTTL  time.Duration
}

type linkEntry struct {
	payload   LinkPayload
	expiresAt time.Time
}

type sessionEntry struct {
	session   Session
	expiresAt time.Time
}

func New(linkTTL, sessionTTL time.Duration) *Service {
	if linkTTL <= 0 {
		linkTTL = 5 * time.Minute
	}
	if sessionTTL <= 0 {
		sessionTTL = 24 * time.Hour
	}
	return &Service{
		links:      make(map[string]linkEntry),
		sessions:   make(map[string]sessionEntry),
		linkTTL:    linkTTL,
		sessionTTL: sessionTTL,
	}
}

func (s *Service) CreateLinkCode(_ context.Context, payload LinkPayload) (string, time.Duration, error) {
	if payload.PlayerID == uuid.Nil {
		return "", 0, fmt.Errorf("player id required")
	}
	code, err := randomCode(8)
	if err != nil {
		return "", 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeLocked(time.Now())
	s.links[strings.ToUpper(code)] = linkEntry{
		payload:   payload,
		expiresAt: time.Now().Add(s.linkTTL),
	}
	return code, s.linkTTL, nil
}

func (s *Service) RedeemLinkCode(_ context.Context, code string) (*Session, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, ErrLinkCodeInvalid
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeLocked(now)
	entry, ok := s.links[code]
	if !ok || now.After(entry.expiresAt) {
		return nil, ErrLinkCodeInvalid
	}
	delete(s.links, code)

	token, err := randomToken()
	if err != nil {
		return nil, err
	}
	session := Session{
		Token:         token,
		PlayerID:      entry.payload.PlayerID,
		MinecraftUUID: entry.payload.MinecraftUUID,
		Username:      entry.payload.Username,
		ExpiresAt:     now.Add(s.sessionTTL),
	}
	s.sessions[token] = sessionEntry{session: session, expiresAt: session.ExpiresAt}
	return &session, nil
}

func (s *Service) ResolveToken(_ context.Context, token string) (*Session, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrSessionInvalid
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.purgeLocked(now)
	entry, ok := s.sessions[token]
	if !ok || now.After(entry.expiresAt) {
		return nil, ErrSessionInvalid
	}
	session := entry.session
	return &session, nil
}

func (s *Service) RevokeToken(_ context.Context, token string) {
	token = strings.TrimSpace(token)
	if token == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

func (s *Service) purgeLocked(now time.Time) {
	for code, entry := range s.links {
		if now.After(entry.expiresAt) {
			delete(s.links, code)
		}
	}
	for token, entry := range s.sessions {
		if now.After(entry.expiresAt) {
			delete(s.sessions, token)
		}
	}
}

func randomCode(n int) (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
