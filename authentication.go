package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret"))

func LoginCorrect(user string, password string) int {
	// 0: Access granted
	// 1: User doesn't exist
	// 2: User exists, but the password is incorrect

	//access := loginCorrect("dartie", "pwd")
	//_ = access

	query := "SELECT id, username, password FROM User WHERE username=?;"
	row := db.QueryRow(query, user)

	var dbid int
	var dbusername, dbpassword string
	err := row.Scan(&dbid, &dbusername, &dbpassword)

	if err != nil {
		return 1 // No rows were returned
	} else {
		if dbpassword == password {
			return 0
		} else {
			return 2
		}
	}
}

func isAuthenticated(handlerFunc func(*gin.Context), c *gin.Context) {
	auth(c)
}

func isLoggedIn(c *gin.Context) bool {
	_, err := c.Cookie("session_token")

	if err == http.ErrNoCookie {
		// user NOT logged in
		return false
	}

	return true
}

// middleware for checking whether the user is admin
func isAdmin(c *gin.Context) {
	userMap := getCurrentUserMap(c)
	var isAdmin bool
	if userMap == nil {
		isAdmin = false
	} else {
		isAdmin = userMap["is_admin"].(bool)
	}

	if !isAdmin {
		c.HTML(http.StatusForbidden, "403-Forbidden.html", nil)
		c.Abort()
		return
	}
	c.Next()
}

// auth middleware checks if logged in by looking at session
func auth(c *gin.Context) {
	nextUrl := c.Request.RequestURI

	session, _ := store.Get(c.Request, "session")
	_, ok := session.Values["user"]
	if !ok {
		//c.HTML(http.StatusForbidden, "login.html", nil)
		c.HTML(http.StatusForbidden, "login.html", gin.H{"nextUrl": nextUrl})
		c.Abort()
		return
	}
	c.Next()
}

/*
	rows := sqlQuery("SELECT id, username, password FROM User;")
	_ = rows
	var id int
	var username, password string
	for rows.Next() {
		// Get values from row.
		err := rows.Scan(&id, &username, &password)
		if err != nil {
			fmt.Print(err)
		}
	}

	defer rows.Close()
*/

// Displays form for login
func getLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func postLoginHandler(c *gin.Context) {
	var user User
	user.Username = c.PostForm("username")
	password := c.PostForm("password")
	defaultRedirect := settingsMap["login_redirect"]
	nextUrl := c.DefaultPostForm("nextUrl", defaultRedirect)
	if nextUrl == "" {
		nextUrl = defaultRedirect
	}

	_ = password
	err := user.getUserByUsername()

	if err != nil {
		fmt.Println("error selecting pswd_hash in db by Username, err:", err)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "check username and password"})
		return
	}
	err = nil //bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))  //TODO : user Hashed PWD instead of clear PWD
	if err == nil {
		session, _ := store.Get(c.Request, "session")
		// session struct has field Values map[interface{}]interface{}
		session.Values["user"] = user
		// save before writing to response/return from handler
		session.Save(c.Request, c.Writer)
		//c.HTML(http.StatusOK, "home.html", gin.H{"username": user.Username})
		c.Redirect(http.StatusMovedPermanently, nextUrl)
		return
	}
	fmt.Println("err:", err)
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "check username and password"})
}

// Logout user by deleting session data
func getLogoutHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	delete(session.Values, "user")
	session.Save(c.Request, c.Writer)
	c.HTML(http.StatusOK, "home.html", gin.H{"message": "Logged out"})
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
