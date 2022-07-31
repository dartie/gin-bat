package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {
	// Login
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.POST("/login", postLoginHandler)

	// Logout
	r.GET("/logout", getLogoutHandler)

	// Register
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "registration.html", nil) })

	protected := r.Group("/").Use(auth)
	unprotected := r.Group("/")
	_ = protected
	_ = unprotected

	protected.GET("/protected", addPerson)
	unprotected.GET("/read", updatePerson)

}
