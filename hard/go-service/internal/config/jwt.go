package config

import (
	"os"
)

// DefaultJWTSecret matches hard/python-service when JWT_SECRET is unset (dev only).
const DefaultJWTSecret = "dev-secret-change-me"

func JWTSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = DefaultJWTSecret
	}
	return []byte(s)
}
