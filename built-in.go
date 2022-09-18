package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/itchyny/timefmt-go"
	"golang.org/x/crypto/bcrypt"
)

func checkErr(err error) {
	if err != nil {
		//log.Panic(err)
		log.Panic(color.Red.Sprintf("%s", err))
	}
}

func checkErrCmd(err error, msg string, errCode int) {
	if err != nil {
		//log.Panic(err)
		log.Fatalf(color.Red.Sprint(msg))
		os.Exit(errCode)
	}
}

func getCheckBoxValue(c *gin.Context, field string) (bool, int) {
	fieldValue := c.PostForm(field)
	var fieldValueInt int   //needed by Postgres
	var fieldValueBool bool // works with other Database
	if fieldValue == "" {
		fieldValueBool = false
		fieldValueInt = 0
	} else {
		fieldValueBool = true
		fieldValueInt = 1
	}
	return fieldValueBool, fieldValueInt
}

// converts all query strings using question marks as paramater "=?" to dollars parameters "=$1, =$2.."
func replaceQuestionMarksWithDollarsInQuery(queryString string) string {
	var newQueryString string
	dollarsCount := 0

	for _, char := range queryString {
		if char == '?' {
			dollarsCount += 1
			newQueryString += "$" + fmt.Sprint(dollarsCount)
		} else {
			newQueryString += string(char)
		}
	}

	return newQueryString
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
		checkErr(getUserError)
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

/* Built-in Additional functions for template */
// return the last element of a slice
func last(x int, a interface{}) bool {
	return x == reflect.ValueOf(a).Len()-1
}

// join elements of a list in gin template
func join(list []string, char string) string {
	return strings.Join(list, char)
}

// make the path from a list of directories + the input element
func makePath(list []string, element string) string {
	newPathList := append(list, []string{element}...)
	newPathString := strings.Join(newPathList, "/")

	return newPathString
}

// unescapes Html text : write the text as it is
func unescapeHtml(s string) template.HTML {
	return template.HTML(s)
}

// make the path from a list of directories + the input element
func setIconColor(filename string) string {
	// "supportedExtensions" defined in "settings.go"
	var style = ""
	ext := filepath.Ext(filename)
	ext = strings.Trim(strings.ToLower(ext), ".")

	if val, ok := supportedExtensions[ext]; ok {
		if val[1] == "" {
			style = ""
		} else {
			style = fmt.Sprintf("color: %s", val[1])
		}
	}

	return style
}

// make the path from a list of directories + the input element
func setFileIcon(filename string) string {
	// "supportedExtensions" defined in "settings.go"
	var classText = "bi-file-earmark"
	ext := filepath.Ext(filename)
	ext = strings.Trim(strings.ToLower(ext), ".")

	if val, ok := supportedExtensions[ext]; ok {
		if val[0] == "" {
			classText = "bi-filetype-" + ext
		} else {
			classText = val[0]
		}
	}

	return classText
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
	_, isAdminInt := getCheckBoxValue(c, "isAdmin")
	picture, _ := c.FormFile("upload_profile_pic")

	var message string
	var status string

	// Generate hashed password from bcrypt
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	checkErr(hashedPasswordErr)

	// Check if new profile already exists
	key := strings.ToUpper(username)
	query := `SELECT username from "Users" where UPPER(username) = $1;`
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
		profileData, err = io.ReadAll(file)
		checkErr(err)
	} else {
		profileData = []uint8{0}
	}

	sqlInsertString := `INSERT INTO "Users" ( 
username, password, first_name, last_name, email, birthday, picture, phone, date_joined, last_login, role, is_admin, active
)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
`
	sqlCommand, err := db.Prepare(sqlInsertString)
	checkErr(err)
	sqlResult, sqlErr := sqlCommand.Exec(username, hashedPassword, firstName, lastName, email, birthday, profileData, phone, nowSqliteFormat(), "", "", isAdminInt, 1)

	if sqlErr == nil {
		recordId, err := sqlResult.LastInsertId() // doesn't work with Postgres.
		if err != nil {
			recordId = 0
		}
		_ = recordId // TODO: print log.
		message = fmt.Sprintf("User \"%s\" has been created successfully", username)
		status = "0"
	} else {
		message = fmt.Sprintf("Issues creating user %s", username)
		status = "1"
	}

	c.HTML(http.StatusOK, "home.html", gin.H{"User": userInfoMap, "Feedback": map[string]string{message: status}, "Url": "/"})
}

// POST Update User
func postUpdateUserHandler(c *gin.Context) {
	updateUserFunc(c, false)
}

func updateUserFunc(c *gin.Context, admin bool) {
	userInfoMap := getCurrentUserMap(c)
	username := c.DefaultPostForm("username", "")
	_, isAdminInt := getCheckBoxValue(c, "isAdmin")
	firstName := c.DefaultPostForm("first_name", "")
	lastName := c.DefaultPostForm("last_name", "")
	phone := c.DefaultPostForm("mobile", "")
	email := c.DefaultPostForm("email", "")
	birthday := dateToDbFormat(c.DefaultPostForm("birthday", ""))
	picture, _ := c.FormFile("upload_profile_pic")

	var sqlUpdateString string

	sqlUpdateString = `UPDATE "Users" SET
	%s
	WHERE
	id = ?;
	`

	queryFields := []string{"username", "first_name", "last_name", "email", "birthday", "phone", "role"}
	var queryParameters []interface{}
	queryParameters = append(queryParameters, username, firstName, lastName, email, birthday, phone, "")

	var profileData []byte
	var sqlErr error
	if picture != nil {
		file, err := picture.Open()
		checkErr(err)
		defer file.Close()
		profileData, err = io.ReadAll(file)
		checkErr(err)
		queryFields = append(queryFields, "picture")
		queryParameters = append(queryParameters, profileData)
	} /*else {
		//profileData = []uint8{0}

	}*/

	if admin {
		queryFields = append(queryFields, "is_admin")
		queryParameters = append(queryParameters, isAdminInt)
	}

	queryFieldsString := strings.Trim(strings.Join(queryFields, "=?, "), ", ") + "=? "
	sqlUpdateString = fmt.Sprintf(sqlUpdateString, queryFieldsString)
	sqlUpdateString = replaceQuestionMarksWithDollarsInQuery(sqlUpdateString)
	sqlCommand, err := db.Prepare(sqlUpdateString)
	checkErr(err)

	// Add id parameter
	queryParameters = append(queryParameters, userInfoMap["id"])

	_, sqlErr = sqlCommand.Exec(queryParameters...)

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

/* Ajax Handlers */
// Validate new User
func validateNewUser(c *gin.Context) {
	// Check if new profile already exists
	username := c.DefaultPostForm("username", "")
	var errors []string

	key := strings.ToUpper(username)
	query := `SELECT username from "Users" where UPPER(username) = $1;`
	row := db.QueryRow(query, key)
	var dbid interface{}
	err := row.Scan(&dbid)
	if err == nil {
		// User already exists
		message := fmt.Sprintf("User %s already exists", username)
		errors = append(errors, message)
	}

	c.JSON(http.StatusOK, gin.H{"errors": errors})
}
