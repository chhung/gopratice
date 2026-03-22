package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireBearerToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handler := RequireBearerToken("secret-token", next)
	req := httptest.NewRequest(http.MethodPost, "/admin/shutdown", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("ServeHTTP() status = %d, want %d", res.Code, http.StatusNoContent)
	}
}

func TestRequireBearerTokenRejectsInvalidToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handler := RequireBearerToken("secret-token", next)
	req := httptest.NewRequest(http.MethodPost, "/admin/shutdown", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("ServeHTTP() status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}
