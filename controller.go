package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPersons(c *gin.Context) {
	//auth(c)
	//c.HTML(http.StatusOK, "index.html", nil)
	c.JSON(http.StatusOK, gin.H{"message": "getPersons Called"})
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "getPersonById " + id + " Called"})
}

func addPerson(c *gin.Context) {
	userMap := getCurrentUserMap(c)
	//fmt.Println(userinfo.Username)

	//c.JSON(http.StatusOK, gin.H{"message": "addPerson Called", "User": userinfo})
	c.HTML(http.StatusOK, "profile.html", gin.H{"User": userMap})
}

func updatePerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updatePerson Called"})
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "deletePerson " + id + " Called"})
}

func options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "options Called"})
}
