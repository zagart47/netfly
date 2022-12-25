package main

import (
	"net/http"
	"netfly/controller"
	"netfly/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/*.html")
	router.Static("/assets", "./assets")

	public := router.Group("/api")
	public.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", gin.H{
			"message": "User registered",
		})
	})
	public.POST("/register", controller.Register)
	public.POST("/login", controller.Login)
	public.GET("/message", controller.ReadMessage)
	public.POST("/message", controller.SendMessage)
	public.GET("/", controller.Profile)
	protected := router.Group("/api/admin")
	protected.Use(middleware.JwtAuthMiddleware())
	router.Run("localhost:8080")
}
