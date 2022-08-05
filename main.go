package main

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

/* Global variables */
var settingsMap map[string]string
var db *sql.DB
var r *gin.Engine

func readSettings() map[string]string {
	/* Read settings */
	var settingsMap map[string]string
	settingsFile := "settings.json"
	if _, err := os.Stat(settingsFile); err == nil {
		settingsBytes, err := ioutil.ReadFile(settingsFile)
		if err != nil {
			log.Fatal(err)
		}
		settingsStr := string(settingsBytes)

		json.Unmarshal([]byte(settingsStr), &settingsMap)

	} else {
		settingsMap["host"] = "0.0.0.0"
		settingsMap["port"] = "8000"
		settingsMap["database_type"] = "sqlite3"
		settingsMap["database_connection"] = "./db.sqlite3"
		settingsMap["logout_redirect"] = "/login"
		settingsMap["login_redirect"] = "/home"
	}

	return settingsMap
}

func main() {
	/* Load settings */
	settingsMap = readSettings()

	// Connect to database
	var dberr error
	db, dberr = sql.Open(settingsMap["database_type"], settingsMap["database_connection"])
	checkErr(dberr)
	// defer close
	defer db.Close()

	/* Initialize Gin */
	r := gin.Default()

	/* Initialize Session */
	// gob.Register(User{}) // Register the User structure
	// store := cookie.NewStore([]byte("snaosnca"))
	// r.Use(sessions.Sessions("SESSIONID", store))

	store.Options.HttpOnly = true // since we are not accessing any cookies w/ JavaScript, set to true
	store.Options.Secure = true   // requires secuire HTTPS connection
	gob.Register(&User{})

	/* Load templates */
	//r.LoadHTMLGlob("templates/**/*")
	var files []string
	filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})

	r.LoadHTMLFiles(files...)

	/* Serve static files */
	r.Static("/static", "./static")

	/* Load all routes from urls.go file */
	routes(r)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	//r.Run(settingsMap["host"] + ":" + settingsMap["port"]) // TODO : restore this line for production
	r.Run("localhost:" + settingsMap["port"]) // only for debug
}
