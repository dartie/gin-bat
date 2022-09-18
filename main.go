package main

import (
	"database/sql"
	"embed"
	"encoding/gob"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

/* Global variables */
var settingsMap map[string]string
var db *sql.DB

//go:embed all:static all:templates *.go db.sqlite3 settings.json go.mod go.sum
var fsProjectFiles embed.FS

func readSettings() map[string]string {
	/* Read settings */
	settingsMap = make(map[string]string)
	settingsFile := "settings.json"
	if _, err := os.Stat(settingsFile); err == nil {
		settingsBytes, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Fatal(err)
		}
		settingsStr := string(settingsBytes)

		// remove json comments
		var settingsStrNoComments string
		for _, settingLine := range strings.Split(settingsStr, "\n") {
			if !strings.HasPrefix(strings.TrimSpace(settingLine), "//") {
				settingsStrNoComments += settingLine + "\n"
			}
		}

		json.Unmarshal([]byte(settingsStrNoComments), &settingsMap)

	} else {
		settingsMap["host"] = "0.0.0.0"
		settingsMap["port"] = "8000"
		settingsMap["database_type"] = "sqlite3"
		settingsMap["database_connection"] = "./db.sqlite3"
		settingsMap["SECRET_KEY"] = RandStringBytesMaskImprSrcUnsafe(60)
		settingsMap["logout_redirect"] = "/login"
		settingsMap["login_redirect"] = "/home"
	}

	return settingsMap
}

func main() {
	/* Load settings */
	settingsMap = readSettings()

	/* Connect to database */
	var dberr error
	db, dberr = sql.Open(settingsMap["database_type"], settingsMap["database_connection"])
	checkErr(dberr)
	// defer close
	defer db.Close()

	/* Init command line */
	initcmd()

	/* Initialize Gin */
	r := gin.Default()

	/* Initialize Session */
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

	/* Load all custom function to pass to the template */
	r.SetFuncMap(template.FuncMap{
		"safe":             unescapeHtml, // Django syntax
		"unescapeHtml":     unescapeHtml,
		"Join":             join,
		"Last":             last,
		"makePath":         makePath,
		"ByteCountDecimal": ByteCountDecimal,
		"SetFileIcon":      setFileIcon,
		"setIconColor":     setIconColor,
	})

	r.LoadHTMLFiles(files...)

	/* Serve static files */
	r.Static("/static", "./static")

	/* Load all routes from urls.go file */
	routes(r)

	/* Run the server */
	if runserverCommand.Invoked {
		// By default it serves on :8080 unless a
		// PORT environment variable was defined.
		//r.Run(settingsMap["host"] + ":" + settingsMap["port"]) // TODO : restore this line for production

		r.Run(*host + ":" + *port)
	}

}
