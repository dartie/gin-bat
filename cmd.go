package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/hellflame/argparse"
	"github.com/howeyc/gopass"
	"github.com/itchyny/timefmt-go"
	"github.com/scylladb/termtables"
	"golang.org/x/crypto/bcrypt"
)

var appName = "MyWebApp"
var Version = "1.0"
var host *string
var port *string

/* Arguments command line */
var parser = argparse.NewParser(appName+" "+Version, ``, nil)
var createUserCommand = parser.AddCommand("create-user", "Create a new user", nil)
var updateUserCommand = parser.AddCommand("update-user", "Update an existing user", nil)
var deleteUserCommand = parser.AddCommand("delete-user", "Delete an existing user", nil)
var resetPasswordCommand = parser.AddCommand("reset-password", "Reset password for an existing user", nil)
var createUserTokenCommand = parser.AddCommand("create-token", "Create a new user", nil)
var displayUserTokenCommand = parser.AddCommand("display-token", "Display token for a given user", nil)
var runserverCommand = parser.AddCommand("runserver", "Run the server", &argparse.ParserConfig{DisableDefaultShowHelp: true})

// Variables for create-user command
var createUserUser *string
var createUserPassword *string
var createUserFirstname *string
var createUserLastname *string
var createUserEmail *string
var createUserPhone *string
var createUserBirthday *string
var createUserPicture *string
var createUserRole *string
var createUserAdmin *bool

func initcmd() {
	// Arguments for createUserCommand
	createUserUser = createUserCommand.String("u", "user", nil)
	createUserPassword = createUserCommand.String("p", "password", nil)
	createUserFirstname = createUserCommand.String("n", "first-name", nil)
	createUserLastname = createUserCommand.String("s", "last-name", nil)
	createUserEmail = createUserCommand.String("e", "email", nil)
	createUserPhone = createUserCommand.String("t", "telephone", nil)
	createUserBirthday = createUserCommand.String("b", "birthday", nil)
	createUserPicture = createUserCommand.String("pic", "picture", nil)
	createUserRole = createUserCommand.String("r", "role", nil)
	createUserAdmin = createUserCommand.Flag("a", "admin", nil)

	_ = createUserUser
	_ = createUserPassword
	_ = createUserFirstname
	_ = createUserLastname
	_ = createUserEmail
	_ = createUserPhone
	_ = createUserBirthday
	_ = createUserPicture
	_ = createUserRole
	_ = createUserAdmin

	// Arguments for updateUserCommand
	updateUserCommand.String("u", "user", nil)
	updateUserCommand.String("p", "password", nil)
	updateUserCommand.String("n", "first-name", nil)
	updateUserCommand.String("s", "last-name", nil)
	updateUserCommand.String("e", "email", nil)
	updateUserCommand.String("t", "telephone", nil)
	//updateUserCommand.String("b", "birthday", &argparse.Option{Help: "Format DD/MM/YYYY"})
	updateUserCommand.String("b", "birthday", &argparse.Option{Help: "Format YYYY-MM-DD"})
	updateUserCommand.String("pic", "picture", nil)
	updateUserCommand.String("r", "role", nil)
	updateUserCommand.Flag("a", "admin", nil)

	// Arguments for deleteUserCommand
	deleteUserCommand.String("u", "user", nil)

	// Arguments for resetPasswordCommand
	resetPasswordCommand.String("u", "user", nil)

	// Arguments for createUserTokenCommand
	createUserTokenCommand.String("u", "user", nil)

	// Arguments for displayUserTokenCommand
	displayUserTokenCommand.String("u", "user", nil)

	// Arguments for runserverCommand
	host = runserverCommand.String("H", "host", &argparse.Option{Default: settingsMap["host"]})
	port = runserverCommand.String("P", "port", &argparse.Option{Default: settingsMap["port"]})

	/* Parse */
	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}

	if createUserCommand.Invoked {
		createUser()
	}
}

func createUser() {
	sqlInsertString := `INSERT INTO User( 
		Id, username, password, first_name, last_name, email, birthday, picture, phone, date_joined, last_login, role, is_admin, active
		)
		VALUES
		(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
	_ = sqlInsertString

	// Check inputs
	if *createUserUser == "" {
		fmt.Print("Enter username:\n>")
		*createUserUser = readStdin()
	}

	// Check if new profile already exists
	key := strings.ToUpper(*createUserUser)
	query := "SELECT username from User where UPPER(username) = ?"
	row := db.QueryRow(query, key)
	var dbid interface{}
	queryErr := row.Scan(&dbid)
	if queryErr == nil {
		log.Fatalf("User %s already exists. Exiting.", *createUserUser)
	}

	if *createUserPassword == "" {
		fmt.Print("Enter password:\n>")
		inputPassword, _ := gopass.GetPasswd()

		fmt.Print("Confirm password:\n>")
		inputConfirmPassword, _ := gopass.GetPasswd()

		if string(inputConfirmPassword) != string(inputPassword) {
			log.Fatal("Passwords don't match.")
		} else {
			*createUserPassword = string(inputPassword)
		}
	}

	// Generate hashed password from bcrypt
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(*createUserPassword), bcrypt.DefaultCost)
	checkErr(hashedPasswordErr)

	if *createUserEmail == "" {
		fmt.Print("Enter email:\n>")
		*createUserEmail = readStdin()
	}

	if *createUserBirthday == "" {
		fmt.Print("Enter birthday (YYYY-MM-DD):\n>")
		//fmt.Print("Enter birthday (DD-MM-YYYY):\n>")
		*createUserBirthday = readStdin()
	}

	if *createUserBirthday != "" {
		// Validate date input
		//re := regexp.MustCompile(`(0?[1-9]|[12][0-9]|3[01])[-/](0?[1-9]|1[012])[-/]((19|20)\d\d)`)
		re := regexp.MustCompile(`((19|20)\d\d)[-/](0?[1-9]|1[012])[-/](0?[1-9]|[12][0-9]|3[01])`)
		if !re.MatchString(*createUserBirthday) {
			log.Fatalf("Invalid date format: %s", *createUserBirthday)
		}

		dateObject, err := dateparse.ParseAny(*createUserBirthday)
		if err != nil {
			log.Fatalf("Invalid date: %s", *createUserBirthday)
		}

		*createUserBirthday = timefmt.Format(dateObject, "%Y-%m-%d")
	}

	if *createUserPicture == "" {
		fmt.Print("Enter picture path:\n>")
		*createUserPicture = readStdin()
	}

	var profileData interface{} //[]byte
	var readPicerr error
	if *createUserPicture != "" {
		//Validate path
		if *createUserPicture != "" {
			if _, err := os.Stat(*createUserPicture); errors.Is(err, os.ErrNotExist) {
				log.Fatalf("File \"%s\" doesn't exists", *createUserPicture)
			}
		}
		userPictureFile, fileErr := os.Open(*createUserPicture)
		checkErr(fileErr)
		defer userPictureFile.Close()
		profileData, readPicerr = io.ReadAll(userPictureFile)
		checkErr(readPicerr)
	} else {
		profileData = nil
	}

	if *createUserPhone == "" {
		fmt.Print("Enter phone:\n>")
		*createUserPhone = readStdin()
	}

	var adminInputString = "NO"
	if !*createUserAdmin {
		var adminInput string
		fmt.Print("Is the user Admin? (Y\\n):\n>")
		adminInput = readStdin()

		if strings.ToUpper(adminInput) == "Y" {
			*createUserAdmin = true
			adminInputString = "YES"
		}
	} else {
		adminInputString = "YES"
	}

	// Print Summary
	table := termtables.CreateTable()
	table.AddHeaders("Info", "DB Field", "Value")
	table.AddRow("Username", "username", *createUserUser)
	table.AddRow("Email", "email", *createUserEmail)
	table.AddRow("Birthday", "birthday", *createUserBirthday)
	table.AddRow("Profile Picture", "picture", *createUserPicture)
	table.AddRow("Phone number", "phone", *createUserPhone)
	table.AddRow("Administrator", "is_admin", adminInputString)
	fmt.Println(table.Render())

	var confirmSummary string
	fmt.Println("Create the user with the above info? [Y\\n]\n>")
	confirmSummary = readStdin()

	if strings.ToUpper(confirmSummary) != "Y" {
		fmt.Println("User aborted")
		os.Exit(0)
	}

	sqlCommand, err := db.Prepare(sqlInsertString)
	checkErr(err)
	sqlResult, sqlErr := sqlCommand.Exec(*createUserUser, hashedPassword, *createUserFirstname, *createUserLastname, *createUserEmail, *createUserBirthday, profileData, *createUserPhone, nowSqliteFormat(), "", "", *createUserAdmin, true)
	var status int64
	var message string
	if sqlErr == nil {
		recordId, err := sqlResult.LastInsertId()
		if err != nil {
			recordId = 0
		}
		message = fmt.Sprintf("User \"%s\" - Id = \"%d\" has been created successfully", *createUserUser, recordId)
		status = recordId
	} else {
		message = fmt.Sprintf("Issues creating user %s : %s", *createUserUser, sqlErr)
		status = 0
	}
	fmt.Println(message)
	os.Exit(int(status))

}
