package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {
	admin := r.Group("/").Use(isAdmin)
	protected := r.Group("/").Use(auth)
	unprotected := r.Group("/")
	_ = admin
	_ = protected
	_ = unprotected

	protected.GET("/protected", addPerson)
	unprotected.GET("/read", updatePerson)

	// Home
	unprotected.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "home.html", gin.H{"User": getCurrentUserMap(c)}) })

	// Login
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{"User": getCurrentUserMap(c)}) })
	r.POST("/login", postLoginHandler)

	// Logout
	r.GET("/logout", getLogoutHandler)

	// Register
	admin.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{"User": getCurrentUserMap(c), "mode": "create"})
	})

	// Profile
	protected.GET("/profile", func(c *gin.Context) {
		userInfoMap := getCurrentUserMap(c)
		c.HTML(http.StatusOK, "profile.html", gin.H{"User": userInfoMap, "mode": "view"})
	})

	// Edit Profile
	protected.GET("/edit-profile", func(c *gin.Context) {
		userInfoMap := getCurrentUserMap(c)
		c.HTML(http.StatusOK, "profile.html", gin.H{"User": userInfoMap, "mode": "edit"})
	})

	// Create Profile GET /register -> POST /create-profile
	r.POST("/create-profile", postCreateUserHandler)

	// Edit Profile
	r.POST("/update-profile", postUpdateUserHandler)

	// Admin
	admin.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", gin.H{"User": getCurrentUserMap(c)})
	})

	/* Ajax routes */
	// Validate form
	admin.POST("/validate-form", validateNewUser)

	/* My routes : Any custom route */

}
