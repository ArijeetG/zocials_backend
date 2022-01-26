package main

import (
	"github.com/gin-gonic/gin"
	h "socials.com/handlers"
)

func main() {
	r := gin.Default()
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})
	r.POST("/auth", func(c *gin.Context) {
		h.SignUp(c)
	})
	r.POST("/signin", func(c *gin.Context) {
		h.Signin(c)
	})
	r.Run(":8000")
}
