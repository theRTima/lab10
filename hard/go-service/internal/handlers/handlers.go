package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type helloURI struct {
	Name string `uri:"name" binding:"required,min=2,max=50"`
}

type profileBody struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=80"`
	Email       string `json:"email" binding:"required,email"`
	Age         int    `json:"age" binding:"required,gte=1,lte=150"`
}

type tokenRequest struct {
	Sub string `json:"sub" binding:"required,min=1,max=128"`
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HelloName(c *gin.Context) {
	var in helloURI
	if err := c.ShouldBindUri(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "hello there " + in.Name,
	})
}

func IssueToken(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req tokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		claims := jwt.RegisteredClaims{
			Subject:   req.Sub,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := t.SignedString(secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": signed,
			"token_type":   "bearer",
			"expires_in":   int(time.Hour / time.Second),
		})
	}
}

func PostProfile(c *gin.Context) {
	var in profileBody
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"display_name": in.DisplayName,
		"email":        in.Email,
		"age":          in.Age,
	})
}
