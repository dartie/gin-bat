package main

import (
	"database/sql"
	"embed"
	"encoding/gob"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

/* Global variables */
var settingsMap map[string]string
var db *sql.DB

//go:embed all:static all:templates *.go db.sqlite3 settings.json go.mod go.sum
var fsProjectFiles embed.FS

func main() {
	/* Load settings */
	settingsFile := "settings.json"
	if _, err := os.Stat(settingsFile); err == nil {
	} else {
		settingsMap["host"] = "0.0.0.0"
		settingsMap["port"] = "8000"
		settingsMap["debug-mode"] = "false"
		settingsMap["database_type"] = "sqlite3"
		settingsMap["database_connection"] = "./db.sqlite3"
		settingsMap["SECRET_KEY"] = RandStringBytesMaskImprSrcUnsafe(60)
		settingsMap["logout_redirect"] = "/login"
		settingsMap["login_redirect"] = "/home"
	}
	settingsMap = readSettings(settingsFile)

	/* Connect to database */
	var dberr error
	db, dberr = sql.Open(settingsMap["database_type"], settingsMap["database_connection"])
	checkErr(dberr)
	// defer close
	defer db.Close()

	/* Init command line */
	initcmd()

	/* Initialize Gin */
	if settingsMap["debug-mode"] == "false" {
		gin.SetMode(gin.ReleaseMode)
	}
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

	startup()

	/* Run the server */
	if runserverCommand.Invoked {
		// By default it serves on :8080 unless a
		// PORT environment variable was defined.
		//r.Run(settingsMap["host"] + ":" + settingsMap["port"]) // TODO : restore this line for production

		r.Run(*host + ":" + *port)
	}
}
