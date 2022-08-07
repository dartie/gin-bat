package main

import (
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
	r.POST("/create-profile", func(c *gin.Context) {
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
		sqlResult, sqlErr := sqlCommand.Exec(username, password, firstName, lastName, email, birthday, profileData, phone, "", "", "", isAdmin, true)

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

		c.HTML(http.StatusOK, "home.html", gin.H{"User": userInfoMap, "mode": "edit", "Feedback": map[string]string{message: status}, "Url": "/"})
	})
}
