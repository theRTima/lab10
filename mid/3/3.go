package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type helloURI struct {
	Name string `uri:"name" binding:"required,min=2,max=50"`
}

type profileBody struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=80"`
	Email       string `json:"email" binding:"required,email"`
	Age         int    `json:"age" binding:"required,gte=1,lte=150"`
}

func newRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/hello/:name", func(c *gin.Context) {
		var in helloURI
		if err := c.ShouldBindUri(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "hello there " + in.Name,
		})
	})

	r.POST("/profile", func(c *gin.Context) {
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
	})

	return r
}

func main() {
	if err := newRouter().Run(":8080"); err != nil {
		panic(err)
	}
}
