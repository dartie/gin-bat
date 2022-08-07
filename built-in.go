package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCheckBoxValue(c *gin.Context, field string) bool {
	fieldValue := c.PostForm(field)
	var fieldValueBool bool
	if fieldValue == "" {
		fieldValueBool = false
	} else {
		fieldValueBool = true
	}
	return fieldValueBool
}

func dateToDbFormat(fieldDate string) string {
	var format = "%s-%s-%s"
	fieldDateSplit := strings.Split(fieldDate, "/")
	var dateDbFormat string
	if len(fieldDateSplit) == 3 {
		dateDbFormat = fmt.Sprintf(format, fieldDateSplit[2], fieldDateSplit[1], fieldDateSplit[0])
	}

	return dateDbFormat
}

/* Get logged user info functions */
func getCurrentUser(c *gin.Context) *User {
	session, _ := store.Get(c.Request, "session")
	userInfo, infoExists := session.Values["user"]

	if userInfo == nil {
		return nil
	}

	userCookie := userInfo.(*User)
	userId := userCookie.Id

	if !infoExists {
		return nil
	}

	var LoggedUser User
	getUserError := LoggedUser.getUserById(userId)

	if getUserError == nil {
		return &LoggedUser
	} else {
		return nil
	}
}

func userInfoToMap(s *User) map[string]interface{} {
	var myMap map[string]interface{}
	data, _ := json.Marshal(s)
	json.Unmarshal(data, &myMap)

	return myMap
}

func getCurrentUserMap(c *gin.Context) map[string]interface{} {
	userinfo := getCurrentUser(c)
	userMap := userInfoToMap(userinfo)

	return userMap
}
