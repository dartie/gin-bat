package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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
	unprotected.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "home.html", nil) })

	// Login
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.POST("/login", postLoginHandler)

	// Logout
	r.GET("/logout", getLogoutHandler)

	// Register
	admin.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "profile.html", gin.H{"mode": "create"}) })

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
	r.POST("/create-profile", func(c *gin.Context) {
		userInfoMap := getCurrentUserMap(c)
		username := c.DefaultPostForm("username", "")
		password := c.DefaultPostForm("password", "")
		firstName := c.DefaultPostForm("first_name", "")
		lastName := c.DefaultPostForm("last_name", "")
		phone := c.DefaultPostForm("mobile", "")
		email := c.DefaultPostForm("email", "")
		birthday := c.DefaultPostForm("birthday", "")
		picture, _ := c.FormFile("upload_profile_pic")

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

		sqlInsertString := fmt.Sprintf(`INSERT INTO User( 
Id, username, password, first_name, last_name, email, birthday, picture, phone, date_joined, last_login, role, is_admin, active
)
VALUES
(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
		sqlCommand, err := db.Prepare(sqlInsertString)
		checkErr(err)
		sqlResult, sqlErr := sqlCommand.Exec(username, password, firstName, lastName, email, birthday, profileData, phone, "", "", "", true, true)

		var message string
		var status string
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

		c.HTML(http.StatusOK, "home.html", gin.H{"User": userInfoMap, "Feedback": map[string]string{message: status}})
	})

	// Edit Profile
	r.POST("/update-profile", func(c *gin.Context) {
		userInfoMap := getCurrentUserMap(c)
		username := c.DefaultPostForm("username", "")
		password := c.DefaultPostForm("password", "")
		firstName := c.DefaultPostForm("first_name", "")
		lastName := c.DefaultPostForm("last_name", "")
		phone := c.DefaultPostForm("mobile", "")
		email := c.DefaultPostForm("email", "")
		birthday := c.DefaultPostForm("birthday", "")
		picture, _ := c.FormFile("upload_profile_pic")

		var profileData []byte
		if picture != nil {
			file, err := picture.Open()
			checkErr(err)
			defer file.Close()
			profileData, err = ioutil.ReadAll(file)
			checkErr(err)
		}

		sqlInsertString := fmt.Sprintf(`INSERT INTO User( 
Id, username, password, first_name, last_name, email, birthday, picture, phone, date_joined, last_login, role, is_admin, active
)
VALUES
(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
		sqlUpdateString := fmt.Sprintf(`UPDATE User SET
username=?, password=?, first_name=?, last_name=?, email=?, birthday=?, picture=?, phone=?, role=?, is_admin=?, active=?
WHERE
id = ?)
`)

		var res sql.Result
		var sqlErr error
		mode := "create"
		if mode == "create" {
			sqlCommand, err := db.Prepare(sqlInsertString)
			checkErr(err)
			res, sqlErr = sqlCommand.Exec(username, password, firstName, lastName, email, birthday, profileData, phone, "", "", "", true, true)
		} else if mode == "edit" {
			sqlCommand, err := db.Prepare(sqlUpdateString)
			checkErr(err)
			res, sqlErr = sqlCommand.Exec(username, password, firstName, lastName, email, birthday, profileData, phone, "", true, true, 1) // TODO: replace 1 with correct 1 or use a different function.
		}
		_ = res
		_ = sqlErr

		c.HTML(http.StatusOK, "profile.html", gin.H{"User": userInfoMap, "mode": "edit"})
	})
}
