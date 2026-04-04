package httpserver

import (
	"lab10/hard/go-service/internal/auth"
	"lab10/hard/go-service/internal/config"
	"lab10/hard/go-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	secret := config.JWTSecret()

	r.GET("/ping", handlers.Ping)
	r.GET("/hello/:name", handlers.HelloName)
	r.POST("/auth/token", handlers.IssueToken(secret))

	protected := r.Group("")
	protected.Use(auth.BearerJWT(secret))
	protected.POST("/profile", handlers.PostProfile)

	return r
}
