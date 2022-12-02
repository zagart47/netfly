package main

import (
	"netfly/controller"
	"netfly/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	public := router.Group("/api")
	public.POST("/register", controller.Register)
	public.POST("/login", controller.Login)

	protected := router.Group("/api/admin")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("/user", controller.CurrentUser)

	router.Run(":8080")
}
