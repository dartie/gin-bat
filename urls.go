package main

import (
	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {
	r.GET("/login", getLoginHandler)
	r.POST("/login", postLoginHandler)
	r.GET("/logout", getLogoutHandler)

	protected := r.Group("/").Use(auth)
	unprotected := r.Group("/")
	_ = protected
	_ = unprotected

	protected.GET("/protected", addPerson)
	unprotected.GET("/read", updatePerson)

}
