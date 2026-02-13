package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MrTomatePNG/projeto-m/internal/auth"
)

func TestRequireAuth_MissingAuthorizationHeader(t *testing.T) {
	jwtm, err := auth.NewJWTManager("secret", time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	called := false
	h := RequireAuth(jwtm)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if called {
		t.Fatalf("expected next handler not to be called")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_InvalidAuthorizationHeader(t *testing.T) {
	jwtm, err := auth.NewJWTManager("secret", time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	called := false
	h := RequireAuth(jwtm)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Token abc")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if called {
		t.Fatalf("expected next handler not to be called")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_ValidTokenSetsUserIDInContext(t *testing.T) {
	jwtm, err := auth.NewJWTManager("secret", time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tok, err := jwtm.Generate(123)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	called := false
	h := RequireAuth(jwtm)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		uid, ok := UserIDFromContext(r.Context())
		if !ok {
			t.Fatalf("expected user id in context")
		}
		if uid != 123 {
			t.Fatalf("expected user id %d, got %d", 123, uid)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if !called {
		t.Fatalf("expected next handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
