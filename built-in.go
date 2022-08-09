package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itchyny/timefmt-go"
	"golang.org/x/crypto/bcrypt"
)

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
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
	if userinfo == nil {
		// Error occurred (user disconnected?)
		logoutUser(c)

		// Options: return nil which raise a Panic Error, or redirect to another view (logout/login).
		return nil
	}
	userMap := userInfoToMap(userinfo)

	return userMap
}

func nowSqliteFormat() string {
	t := time.Now()
	str := timefmt.Format(t, "%Y-%m-%d %H:%M:%S.%f") // YYYY-MM-DD HH:MM:SS.SSS

	return str
}

/* Built-in Handlers */

// POST Create Profile
func postCreateUserHandler(c *gin.Context) {
	userInfoMap := getCurrentUserMap(c)
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	firstName := c.DefaultPostForm("first_name", "")
	lastName := c.DefaultPostForm("last_name", "")
	phone := c.DefaultPostForm("mobile", "")
	email := c.DefaultPostForm("email", "")
	birthday := c.DefaultPostForm("birthday", "")
	isAdmin := getCheckBoxValue(c, "isAdmin")
	picture, _ := c.FormFile("upload_profile_pic")

	var message string
	var status string

	// Generate hashed password from bcrypt
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkErr(hashedPasswordErr)

	// Check if new profile already exists
	key := username
	query := "SELECT username from User where username = ?"
	row := db.QueryRow(query, key)
	var dbid interface{}
	err := row.Scan(&dbid)
	if err == nil {
		message = fmt.Sprintf("User %s already exists", username)
		status = "2"
		c.HTML(http.StatusOK, "home.html", gin.H{"User": userInfoMap, "Feedback": map[string]string{message: status}, "Url": "/"})
		return
	}

	var profileData []byte
	if picture != nil {
		file, err := picture.Open()
		checkErr(err)
		defer file.Close()
		profileData, err = ioutil.ReadAll(file)
		checkErr(err)
	} else {
		profileData = []uint8{0}
	}

	sqlInsertString := `INSERT INTO User( 
Id, username, password, first_name, last_name, email, birthday, picture, phone, date_joined, last_login, role, is_admin, active
)
VALUES
(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	sqlCommand, err := db.Prepare(sqlInsertString)
	checkErr(err)
	sqlResult, sqlErr := sqlCommand.Exec(username, hashedPassword, firstName, lastName, email, birthday, profileData, phone, nowSqliteFormat(), "", "", isAdmin, true)

	if sqlErr == nil {
		recordId, err := sqlResult.LastInsertId()
		if err != nil {
			recordId = 0
		}
		message = fmt.Sprintf("User \"%s\" - Id = \"%d\" has been created successfully", username, recordId)
		status = "0"
	} else {
		message = fmt.Sprintf("Issues creating user %s", username)
		status = "1"
	}

	c.HTML(http.StatusOK, "home.html", gin.H{"User": userInfoMap, "Feedback": map[string]string{message: status}, "Url": "/"})
}

// POST Update User
func postUpdateUserHandler(c *gin.Context) {
	userInfoMap := getCurrentUserMap(c)
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	firstName := c.DefaultPostForm("first_name", "")
	lastName := c.DefaultPostForm("last_name", "")
	phone := c.DefaultPostForm("mobile", "")
	email := c.DefaultPostForm("email", "")
	birthday := dateToDbFormat(c.DefaultPostForm("birthday", ""))
	isAdmin := getCheckBoxValue(c, "isAdmin")
	picture, _ := c.FormFile("upload_profile_pic")

	var profileData []byte
	var sqlErr error
	if picture != nil {
		file, err := picture.Open()
		checkErr(err)
		defer file.Close()
		profileData, err = ioutil.ReadAll(file)
		checkErr(err)
		sqlUpdateString := `UPDATE User SET
		username=?, password=?, first_name=?, last_name=?, email=?, birthday=?, phone=?, role=?, is_admin=?, picture=?
		WHERE
		id = ?;
		`
		sqlCommand, err := db.Prepare(sqlUpdateString)
		checkErr(err)
		_, sqlErr = sqlCommand.Exec(username, password, firstName, lastName, email, birthday, phone, "", isAdmin, profileData, userInfoMap["id"])
	} else {
		//profileData = []uint8{0}
		sqlUpdateString := `UPDATE User SET
		username=?, password=?, first_name=?, last_name=?, email=?, birthday=?, phone=?, role=?, is_admin=?
		WHERE
		id = ?;
		`
		sqlCommand, err := db.Prepare(sqlUpdateString)
		checkErr(err)
		_, sqlErr = sqlCommand.Exec(username, password, firstName, lastName, email, birthday, phone, "", isAdmin, userInfoMap["id"])
	}

	var message string
	var status string
	if sqlErr == nil {
		message = fmt.Sprintf("User \"%s\" - Id = \"%.0f\" has been updated successfully", username, userInfoMap["id"])
		status = "0"
	} else {
		message = fmt.Sprintf("Issues creating user %s", username)
		status = "1"
	}

	// refresh date as they are changed. This is useful for updating the profile picture in the topbar if present.
	userInfoMap = getCurrentUserMap(c)

	c.HTML(http.StatusOK, "home.html", gin.H{"User": userInfoMap, "mode": "edit", "Feedback": map[string]string{message: status}, "Url": "/"})
}
