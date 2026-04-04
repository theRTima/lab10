package httpserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"lab10/hard/go-service/internal/config"
	"lab10/hard/go-service/internal/httpserver"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestAuthToken_GeneratesParsableJWT(t *testing.T) {
	r := httpserver.NewRouter()
	body, _ := json.Marshal(map[string]string{"sub": "user-1"})
	req := httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var out struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.TokenType != "bearer" || out.AccessToken == "" || out.ExpiresIn <= 0 {
		t.Fatalf("unexpected token response: %+v", out)
	}

	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(out.AccessToken, claims, func(tok *jwt.Token) (any, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return config.JWTSecret(), nil
	})
	if err != nil {
		t.Fatalf("parse jwt: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Fatalf("sub: want user-1, got %q", claims.Subject)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Fatal("expected future exp")
	}
}

func TestProfile_WithoutToken_Unauthorized(t *testing.T) {
	r := httpserver.NewRouter()
	body := []byte(`{"display_name":"Tim","email":"t@example.com","age":20}`)
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestProfile_WithToken_InvalidBody_BadRequest(t *testing.T) {
	r := httpserver.NewRouter()
	token := mustIssueToken(t, r, "u42")

	invalid := []byte(`{"display_name":"Tim","email":"not-an-email","age":20}`)
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(invalid))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d, body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func mustIssueToken(t *testing.T, h http.Handler, sub string) string {
	t.Helper()
	body, err := json.Marshal(map[string]string{"sub": sub})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("issue token: status %d %s", rec.Code, rec.Body.String())
	}
	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out.AccessToken
}
