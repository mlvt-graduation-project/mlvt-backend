package main

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.String(200, "Welcome to Home")
}
