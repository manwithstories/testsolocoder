package service

import (
	"errors"
	"sync"
	"time"
)

type VerificationCode struct {
	Code      string
	ExpiresAt time.Time
}

type VerificationCodeStore struct {
	codes map[string]*VerificationCode
	mu    sync.RWMutex
}

var (
	verificationStore *VerificationCodeStore
	once              sync.Once
)

func GetVerificationStore() *VerificationCodeStore {
	once.Do(func() {
		verificationStore = &VerificationCodeStore{
			codes: make(map[string]*VerificationCode),
		}
		go verificationStore.startCleanup()
	})
	return verificationStore
}

func (s *VerificationCodeStore) Store(email, code string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[email] = &VerificationCode{
		Code:      code,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (s *VerificationCodeStore) Verify(email, code string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored, exists := s.codes[email]
	if !exists {
		return errors.New("verification code not found or expired")
	}

	if time.Now().After(stored.ExpiresAt) {
		delete(s.codes, email)
		return errors.New("verification code has expired")
	}

	if stored.Code != code {
		return errors.New("invalid verification code")
	}

	delete(s.codes, email)
	return nil
}

func (s *VerificationCodeStore) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for email, code := range s.codes {
			if now.After(code.ExpiresAt) {
				delete(s.codes, email)
			}
		}
		s.mu.Unlock()
	}
}
