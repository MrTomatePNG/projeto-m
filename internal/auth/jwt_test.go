package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndVerify(t *testing.T) {
	m, err := NewJWTManager("secret", time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	userID := int64(42)
	tok, err := m.Generate(userID)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}
	if tok == "" {
		t.Fatalf("expected token to be non-empty")
	}

	got, err := m.Verify(tok)
	if err != nil {
		t.Fatalf("expected no error verifying token, got %v", err)
	}
	if got != userID {
		t.Fatalf("expected userID %d, got %d", userID, got)
	}
}

func TestJWTManager_Verify_InvalidToken(t *testing.T) {
	m, err := NewJWTManager("secret", time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = m.Verify("not-a-token")
	if err == nil {
		t.Fatalf("expected error verifying invalid token")
	}
}

func TestJWTManager_NewJWTManager_EmptySecret(t *testing.T) {
	_, err := NewJWTManager("", time.Hour)
	if err == nil {
		t.Fatalf("expected error when secret is empty")
	}
}
