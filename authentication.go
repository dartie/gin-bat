package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gookit/color"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte(settingsMap["SECRET_KEY"])
var store = sessions.NewCookieStore(secret)

func LoginCorrect(user string, password string) int {
	// 0: Access granted
	// 1: User doesn't exist
	// 2: User exists, but the password is incorrect

	//access := loginCorrect("dartie", "pwd")
	//_ = access

	query := `SELECT id, username, password FROM "Users" WHERE active=1 AND username=$1;`
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

// auth middleware checks if token is correct and the user allowed
func authToken(c *gin.Context) {
	// Get Token from request
	tokenRequest := c.Request.Header["Authorization"]

	var token string
	if len(tokenRequest) > 0 {
		token = strings.TrimPrefix(tokenRequest[0], "Bearer ")
	}

	var requestorId float64
	claims, claimsErr := GetClaimsFromToken(token)
	if claimsErr != nil {
		errorMessage := fmt.Sprintf("Invalid token: %s\n\n%s\n", token, claimsErr)
		color.Red.Println(errorMessage)
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
		c.Abort()
		return
	} else {
		requestorId = claims["UserInfo"].(map[string]interface{})["id"].(float64)

		// Check the token in the DB
		query := `SELECT key, user_id FROM "AuthToken" WHERE key=$1;`
		row := db.QueryRow(query, token)

		var dbkey string
		var dbuserid float64
		queryErr := row.Scan(&dbkey, &dbuserid)
		if queryErr == nil {
			// User in the db
			if dbuserid != requestorId {
				c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": "Token user and DB User don't match!"})
				c.Abort()
				return
			}
		} else {
			// User missing from db
			c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": "Unauthorized Token!"})
			c.Abort()
			return
		}
	}
}

// post Login handler: checks whether the user is granted and store the cookies for the session
func postLoginHandler(c *gin.Context) {
	var user User
	user.Username = c.PostForm("username")
	password := c.PostForm("password")
	defaultRedirect := settingsMap["login_redirect"]
	nextUrl := c.DefaultPostForm("nextUrl", defaultRedirect)
	if nextUrl == "" {
		nextUrl = defaultRedirect
	}

	// Get user information from the database
	err := user.getUserByUsername()

	if err != nil {
		fmt.Println("error selecting pswd_hash in db by Username, err:", err)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Feedback": map[string]string{"check username and password": "1"}, "Url": "/"})
		return
	}

	// Compare hashed password stored with the hashed password input
	isPasswordWrong := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if isPasswordWrong == nil {
		session, _ := store.Get(c.Request, "session")
		// session struct has field Values map[interface{}]interface{}
		session.Values["user"] = user

		// save before writing to response/return from handler
		session.Save(c.Request, c.Writer)

		// Update the database with last_login field
		sqlInsertString := `UPDATE "Users" SET last_login=$1 WHERE id=$2`
		sqlCommand, err := db.Prepare(sqlInsertString)
		checkErr(err)
		_, sqlErr := sqlCommand.Exec(nowSqliteFormat(), user.Id)
		if sqlErr != nil {
			log.Panic(sqlErr)
		}

		//c.HTML(http.StatusOK, "home.html", gin.H{"username": user.Username})
		c.Redirect(http.StatusMovedPermanently, nextUrl)
		return
	}
	fmt.Println("err:", err)
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"message": "check username and password"})
}

func logoutUser(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	delete(session.Values, "user")
	session.Save(c.Request, c.Writer)
}

// Logout user by deleting session data
func getLogoutHandler(c *gin.Context) {
	logoutUser(c)
	c.HTML(http.StatusOK, "home.html", gin.H{"Feedback": map[string]string{"Logged out": "2"}, "Url": "/"})
}

// Function for creating a auth user token
type MyJWTClaims struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
}

// Function for creating a auth user token
func CreateToken(sub string, userInfo interface{}, expiration time.Time) (string, error) {
	// Get the token instance with the Signing method
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// Add your claims
	token.Claims = &MyJWTClaims{
		&jwt.RegisteredClaims{
			// Set the exp and sub claims. sub is usually the userID
			ExpiresAt: jwt.NewNumericDate(expiration),
			Subject:   sub,
		},
		userInfo,
	}
	// Sign the token with your secret key
	val, err := token.SignedString(secret)

	if err != nil {
		// On error return the error
		return "", err
	}
	// On success return the token string
	return val, nil
}

func GetClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
