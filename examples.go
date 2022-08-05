package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func examples(r *gin.Engine) {
	admin := r.Group("/").Use(isAdmin)
	protected := r.Group("/").Use(auth)
	unprotected := r.Group("/")
	_ = admin
	_ = protected
	_ = unprotected

	// raise feedback alert in the template page. This is feasible as topbar.html contains the code for it.
	unprotected.GET("/Feedback", func(c *gin.Context) {
		Feedback := map[string]string{"Message": "1", "Message2": "2", "Info": "-1", "Well done!": "0"}
		c.HTML(http.StatusOK, "home.html", gin.H{"Feedback": Feedback})
	})
}
