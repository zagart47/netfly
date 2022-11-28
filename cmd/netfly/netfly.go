package main

import (
	"net/http"
	"netfly/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome!")
	})

	public := router.Group("/api")
	public.POST("/register", controller.Register)

	router.Run(":8080")
}
