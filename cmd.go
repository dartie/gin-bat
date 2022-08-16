package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gookit/color"
	"github.com/hellflame/argparse"
	"github.com/howeyc/gopass"
	"github.com/itchyny/timefmt-go"
	"github.com/scylladb/termtables"
	str2duration "github.com/xhit/go-str2duration/v2"
	"golang.org/x/crypto/bcrypt"
)

var appName = "MyWebApp"
var Version = "1.0"

/* Arguments command line */
var parser = argparse.NewParser(appName+" "+Version, ``, nil)
var createUserCommand = parser.AddCommand("create-user", "Create a new user", nil)
var updateUserCommand = parser.AddCommand("update-user", "Update an existing user", nil)
var deleteUserCommand = parser.AddCommand("delete-user", "Delete an existing user", nil)
var changePasswordCommand = parser.AddCommand("change-password", "Change password for an existing user", nil)
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

// Variables for update-user command
var updateUserUser *string
var updateUserPassword *string
var updateUserFirstname *string
var updateUserLastname *string
var updateUserEmail *string
var updateUserPhone *string
var updateUserBirthday *string
var updateUserPicture *string
var updateUserRole *string
var updateUserAdmin *bool

// Variables for delete-user command
var deleteUserUser *string

// Variables for reset-password command
var changePasswordUser *string

// Variables for create-token command
var createUserTokenUser *string
var createUserTokenExpiration *string
var createUserTokenStrictExpiration *bool

// Variables for display-token command
var displayUserTokenUser *string

// Variables for runserver command
var host *string
var port *string

func initcmd() {
	// Arguments for createUserCommand
	createUserUser = createUserCommand.String("u", "user", nil)
	createUserPassword = createUserCommand.String("p", "password", nil)
	createUserFirstname = createUserCommand.String("n", "first-name", nil)
	createUserLastname = createUserCommand.String("s", "last-name", nil)
	createUserEmail = createUserCommand.String("e", "email", nil)
	createUserPhone = createUserCommand.String("t", "telephone", nil)
	createUserBirthday = createUserCommand.String("b", "birthday", &argparse.Option{Help: "Format YYYY-MM-DD"})
	createUserPicture = createUserCommand.String("pic", "picture", nil)
	createUserRole = createUserCommand.String("r", "role", nil)
	createUserAdmin = createUserCommand.Flag("a", "admin", nil)

	// Arguments for updateUserCommand
	updateUserUser = updateUserCommand.String("u", "user", nil)
	updateUserPassword = updateUserCommand.String("p", "password", nil)
	updateUserFirstname = updateUserCommand.String("n", "first-name", nil)
	updateUserLastname = updateUserCommand.String("s", "last-name", nil)
	updateUserEmail = updateUserCommand.String("e", "email", nil)
	updateUserPhone = updateUserCommand.String("t", "telephone", nil)
	updateUserBirthday = updateUserCommand.String("b", "birthday", &argparse.Option{Help: "Format YYYY-MM-DD"})
	updateUserPicture = updateUserCommand.String("pic", "picture", nil)
	updateUserRole = updateUserCommand.String("r", "role", nil)
	updateUserAdmin = updateUserCommand.Flag("a", "admin", nil)

	// Arguments for deleteUserCommand
	deleteUserUser = deleteUserCommand.String("u", "user", nil)

	// Arguments for changePasswordCommand
	changePasswordUser = changePasswordCommand.String("u", "user", nil)

	// Arguments for createUserTokenCommand
	createUserTokenUser = createUserTokenCommand.String("u", "user", nil)
	createUserTokenExpiration = createUserTokenCommand.String("e", "expiration", &argparse.Option{Default: "24h"})
	createUserTokenStrictExpiration = createUserTokenCommand.Flag("s", "strict-expiration", &argparse.Option{Help: "If enabled, the token expires at the exact time specified by the user and not at the midnight."})

	// Arguments for displayUserTokenCommand
	displayUserTokenUser = displayUserTokenCommand.String("u", "user", nil)

	// Arguments for runserverCommand
	host = runserverCommand.String("H", "host", &argparse.Option{Default: settingsMap["host"]})
	port = runserverCommand.String("P", "port", &argparse.Option{Default: settingsMap["port"]})

	/* Parse */
	if e := parser.Parse(nil); e != nil {
		color.Red.Println(e.Error())
		return
	}

	if createUserCommand.Invoked {
		createUser()
	}

	if deleteUserCommand.Invoked {
		deleteUser()
	}

	if changePasswordCommand.Invoked {
		changePassword()
	}

	if createUserTokenCommand.Invoked {
		createUserToken()
	}

	if displayUserTokenCommand.Invoked {
		displayUserToken()
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
		color.Cyan.Print("Enter username:\n>")
		*createUserUser = readStdin()
	}

	// Check if new profile already exists
	key := strings.ToUpper(*createUserUser)
	query := "SELECT username from User where UPPER(username) = ?"
	row := db.QueryRow(query, key)
	var dbid interface{}
	queryErr := row.Scan(&dbid)
	if queryErr == nil {
		log.Fatalf(color.Red.Sprintf("User %s already exists. Exiting.", *createUserUser))
	}

	if *createUserPassword == "" {
		color.Cyan.Print("Enter password:\n>")
		inputPassword, _ := gopass.GetPasswd()

		color.Cyan.Print("Confirm password:\n>")
		inputConfirmPassword, _ := gopass.GetPasswd()

		if string(inputConfirmPassword) != string(inputPassword) {
			log.Fatal(color.Red.Sprintf("Passwords don't match."))
		} else {
			*createUserPassword = string(inputPassword)
		}
	}

	// Generate hashed password from bcrypt
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(*createUserPassword), bcrypt.DefaultCost)
	checkErr(hashedPasswordErr)

	if *createUserFirstname == "" {
		color.Cyan.Print("Enter First Name:\n>")
		*createUserFirstname = readStdin()
	}

	if *createUserLastname == "" {
		color.Cyan.Print("Enter Last Name:\n>")
		*createUserLastname = readStdin()
	}

	if *createUserEmail == "" {
		color.Cyan.Print("Enter email:\n>")
		*createUserEmail = readStdin()
	}

	if *createUserBirthday == "" {
		color.Cyan.Print("Enter birthday (YYYY-MM-DD):\n>")
		//color.Cyan.Print("Enter birthday (DD-MM-YYYY):\n>")
		*createUserBirthday = readStdin()
	}

	if *createUserBirthday != "" {
		// Validate date input
		//re := regexp.MustCompile(`(0?[1-9]|[12][0-9]|3[01])[-/](0?[1-9]|1[012])[-/]((19|20)\d\d)`)
		re := regexp.MustCompile(`((19|20)\d\d)[-/](0?[1-9]|1[012])[-/](0?[1-9]|[12][0-9]|3[01])`)
		if !re.MatchString(*createUserBirthday) {
			log.Fatalf(color.Red.Sprintf("Invalid date format: %s", *createUserBirthday))
		}

		dateObject, err := dateparse.ParseAny(*createUserBirthday)
		if err != nil {
			log.Fatalf(color.Red.Sprintf("Invalid date: %s", *createUserBirthday))
		}

		*createUserBirthday = timefmt.Format(dateObject, "%Y-%m-%d")
	}

	if *createUserPicture == "" {
		color.Cyan.Print("Enter picture path:\n>")
		*createUserPicture = readStdin()
	}

	var profileData interface{} //[]byte
	var readPicerr error
	if *createUserPicture != "" {
		//Validate path
		if *createUserPicture != "" {
			if _, err := os.Stat(*createUserPicture); errors.Is(err, os.ErrNotExist) {
				log.Fatal(color.Red.Sprintf("File \"%s\" doesn't exists", *createUserPicture))
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
		color.Cyan.Print("Enter phone:\n>")
		*createUserPhone = readStdin()
	}

	var adminInputString = "NO"
	if !*createUserAdmin {
		var adminInput string
		color.Cyan.Print("Is the user Admin? (Y\\n):\n>")
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
	table.AddRow("First Name", "first_name", *createUserFirstname)
	table.AddRow("Last Name", "last_name", *createUserLastname)
	table.AddRow("Email", "email", *createUserEmail)
	table.AddRow("Birthday", "birthday", *createUserBirthday)
	table.AddRow("Profile Picture", "picture", *createUserPicture)
	table.AddRow("Phone number", "phone", *createUserPhone)
	table.AddRow("Administrator", "is_admin", adminInputString)
	color.Cyan.Println(table.Render())

	var confirmSummary string
	color.Yellow.Println("Create the user with the above info? [Y\\n]\n>")
	confirmSummary = readStdin()

	if strings.ToUpper(confirmSummary) != "Y" {
		color.Red.Println("User aborted")
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
		color.Green.Println(message)

	} else {
		message = fmt.Sprintf("Issues creating user %s : %s", *createUserUser, sqlErr)
		status = 0
		color.Red.Println(message)

	}
	os.Exit(int(status))
}

func deleteUser() {
	sqlDeleteString := `DELETE FROM User WHERE username=?`

	if *deleteUserUser == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the user delete (username)\n>")
		*deleteUserUser = readStdin()
	}

	// Check if the user exists
	query := "SELECT Id, username, first_name, last_name, email, birthday, phone, date_joined, last_login, role, is_admin, active FROM User WHERE username=?;"
	row := db.QueryRow(query, *deleteUserUser)

	var dbid int
	var dbusername, dbfirstname, dblastname, dbemail, dbbirthday, dbphone, dbdatejoined, dblastlogin, dbrole, dbisadmin, dbactive string
	selectErr := row.Scan(&dbid, &dbusername, &dbfirstname, &dblastname, &dbemail, &dbbirthday, &dbphone, &dbdatejoined, &dblastlogin, &dbrole, &dbisadmin, &dbactive)

	if selectErr == nil {
		// Print Summary
		table := termtables.CreateTable()
		table.AddHeaders("Info", "DB Field", "Value")
		table.AddRow("Username", "username", dbusername)
		table.AddRow("First Name", "first_name", dbfirstname)
		table.AddRow("Last Name", "last_name", dblastname)
		table.AddRow("Email", "email", dbemail)
		table.AddRow("Birthday", "birthday", dbbirthday)
		table.AddRow("Phone number", "phone", dbphone)
		table.AddRow("Date joined", "date_joined", dbdatejoined)
		table.AddRow("Last login", "last_login", dblastlogin)
		table.AddRow("Role", "role", dbrole)
		table.AddRow("Administrator", "is_admin", dbisadmin)
		color.Cyan.Println(table.Render())
	} else {
		log.Fatalf(color.Red.Sprintf("User %s does not exist.", *deleteUserUser))
	}

	// Ask confirm
	var confirm string
	color.Yellow.Printf("Are you sure you want to delete \"%s\" user? [Y\\n]\n>", *deleteUserUser)
	confirm = readStdin()

	if strings.ToUpper(confirm) == "Y" {
		sqlCommand, err := db.Prepare(sqlDeleteString)
		checkErr(err)
		sqlResult, sqlErr := sqlCommand.Exec(*deleteUserUser)
		checkErr(sqlErr)
		recordId, err := sqlResult.LastInsertId()
		if err != nil {
			recordId = 0
		}
		_ = recordId
		color.Green.Printf("User \"%s\" removed\n", *deleteUserUser)
		os.Exit(0)
	} else {
		color.Red.Println("User aborted")
		os.Exit(0)
	}
}

func changePassword() {
	if *changePasswordUser == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the user to change the password (username)\n>")
		*changePasswordUser = readStdin()
	}

	// Check if the user exists
	query := "SELECT Id, username, first_name, last_name, email, birthday, phone, date_joined, last_login, role, is_admin, active FROM User WHERE username=?;"
	row := db.QueryRow(query, *changePasswordUser)

	var dbid int
	var dbusername, dbfirstname, dblastname, dbemail, dbbirthday, dbphone, dbdatejoined, dblastlogin, dbrole, dbisadmin, dbactive string
	selectErr := row.Scan(&dbid, &dbusername, &dbfirstname, &dblastname, &dbemail, &dbbirthday, &dbphone, &dbdatejoined, &dblastlogin, &dbrole, &dbisadmin, &dbactive)

	if selectErr == nil {
		// Print Summary
		table := termtables.CreateTable()
		table.AddHeaders("Info", "DB Field", "Value")
		table.AddRow("Username", "username", dbusername)
		table.AddRow("First Name", "first_name", dbfirstname)
		table.AddRow("Last Name", "last_name", dblastname)
		table.AddRow("Email", "email", dbemail)
		table.AddRow("Birthday", "birthday", dbbirthday)
		table.AddRow("Phone number", "phone", dbphone)
		table.AddRow("Date joined", "date_joined", dbdatejoined)
		table.AddRow("Last login", "last_login", dblastlogin)
		table.AddRow("Role", "role", dbrole)
		table.AddRow("Administrator", "is_admin", dbisadmin)
		color.Cyan.Println(table.Render())
	} else {
		log.Fatalf(color.Red.Sprintf("User \"%s\" does not exist.", *changePasswordUser))
	}

	color.Cyan.Print("Enter password:\n>")
	inputPassword, _ := gopass.GetPasswd()

	color.Cyan.Print("Confirm password:\n>")
	inputConfirmPassword, _ := gopass.GetPasswd()

	if string(inputConfirmPassword) != string(inputPassword) {
		log.Fatal(color.Red.Sprintf("Passwords don't match."))
	}

	// Generate hashed password from bcrypt
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	checkErr(hashedPasswordErr)

	// Update User with new password
	updatePasswordSql := `UPDATE User SET password=? WHERE id=?`

	sqlCommand, err := db.Prepare(updatePasswordSql)
	checkErr(err)
	_, sqlErr := sqlCommand.Exec(hashedPassword, dbid)
	if sqlErr == nil {
		color.Green.Println("\nPassword updated")
		os.Exit(0)
	} else {
		log.Fatal(color.Red.Sprintf("Errors during password update: %s", sqlErr))
	}
}

func createUserToken() {
	var u User
	if *createUserTokenUser == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the user to change the password (username)\n>")
		*createUserTokenUser = readStdin()
	}

	if *createUserTokenExpiration == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the token duration (date DD-MM-YYYY, duration - 24h, 8h - permanent, )\n>")
		*createUserTokenExpiration = readStdin()
	}

	// Check if the user exists
	query := "SELECT Id, username, first_name, last_name, email, birthday, phone, date_joined, last_login, role, is_admin, active FROM User WHERE username=?;"
	row := db.QueryRow(query, *createUserTokenUser)
	selectErr := row.Scan(&u.Id, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Birthday, &u.Phone, &u.DateJoined, &u.LastLogin, &u.Role, &u.IsAdmin, &u.Active)

	if selectErr != nil {
		log.Fatalf(color.Red.Sprintf("User \"%s\" does not exist.", *createUserTokenUser))
	}

	// Create user if your conditions match. Below, all username and passwords are accepted.
	user := &User{
		Id:       u.Id,
		Username: u.Username,
		Password: u.Password,
	}

	// create time object for duration
	var expiration time.Time
	if strings.Contains(*createUserTokenExpiration, "-") || strings.Contains(*createUserTokenExpiration, "/") {
		// date
		var dateParseErr error
		*createUserTokenExpiration = strings.ReplaceAll(*createUserTokenExpiration, "-", "/")
		*createUserTokenExpiration = *createUserTokenExpiration + " 23:59:59"
		expiration, dateParseErr = timefmt.Parse(*createUserTokenExpiration, "%d/%m/%Y %H:%M:%S")
		checkErrCmd(dateParseErr, fmt.Sprintf("%s", dateParseErr), 1)

	} else if strings.ToLower(*createUserTokenExpiration) == "permanent" {
		// permanent
		expiration = time.Now().AddDate(99, 0, 0) // 99 Years

	} else {
		// duration
		d, durationErr := str2duration.ParseDuration(*createUserTokenExpiration) // Supported formats by me: 1m, 1h, 1d
		checkErrCmd(durationErr, fmt.Sprintf("Invalid duration: %s", *createUserTokenExpiration), 1)

		expiration = time.Now().Add(time.Minute * time.Duration(d.Minutes()))

		// If hours is not specified by the user, the token will last until midnight of the expiration day.
		var parseErr error
		if !strings.Contains(strings.ToLower(*createUserTokenExpiration), "h") && !*createUserTokenStrictExpiration {
			expirationString := timefmt.Format(expiration, "%d/%m/%Y 23:59:59")
			expiration, parseErr = timefmt.Parse(expirationString, "%d/%m/%Y %H:%M:%S")
		}

		checkErrCmd(parseErr, fmt.Sprintf("%s", parseErr), 1)
	}

	// Expiration format for DB: YYYY-MM-DD HH:MM:SS.SSS
	expirationDBFormat := timefmt.Format(expiration, "%Y-%m-%d %H:%M:%S.%f")

	// Create token
	tokenString, _ := CreateToken(fmt.Sprintf("%d", u.Id), user, expiration)

	/* Store the token in the DB */
	// Check whether it's and INSERT or UPDATE
	var dbId int
	var sqlInsertString string

	query = "SELECT user_id FROM AuthToken WHERE user_id=?"
	row = db.QueryRow(query, u.Id)
	selectErr = row.Scan(&dbId)
	if selectErr == nil {
		sqlInsertString = `UPDATE AuthToken SET key=?, created=?, expiration=? WHERE user_id=?`
	} else {
		// A record is already present -> UPDATE
		sqlInsertString = `INSERT INTO AuthToken ( 
			key, created, expiration, user_id
			)
			VALUES
			(?, ?, ?, ?)
			`
	}

	// Execute Query
	sqlCommand, err := db.Prepare(sqlInsertString)
	checkErr(err)
	_, sqlErr := sqlCommand.Exec(tokenString, nowSqliteFormat(), expirationDBFormat, u.Id)
	checkErr(sqlErr)

	// Display the token
	fmt.Println() // Blank line
	color.Green.Println(tokenString)
	fmt.Println() // Blank line

	if strings.ToLower(*createUserTokenExpiration) == "permanent" {
		color.Yellow.Println("The above token is permanent (does not expire unless manually revoked).")
	} else {
		color.White.Print("expires in date: ")
		color.Yellow.Println(timefmt.Format(expiration, "%d/%m/%Y (%A %d %B %Y) %H:%M:%S %p"))
	}
	fmt.Println() // Blank line

	os.Exit(0)
}

func displayUserToken() {
	if *displayUserTokenUser == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the user you want to display the token\n>")
		*displayUserTokenUser = readStdin()
	}

	query := `SELECT key, expiration
		FROM AuthToken
		INNER JOIN User
		ON AuthToken.user_id = User.Id
		WHERE User.username=?;
		`

	var dbKey string
	var dbKeyExpiration string
	row := db.QueryRow(query, *displayUserTokenUser)
	selectErr := row.Scan(&dbKey, &dbKeyExpiration)
	if selectErr != nil {
		log.Fatalf(color.Red.Sprintf("User \"%s\" doesn't have tokens", *displayUserTokenUser))
		os.Exit(1)
	}

	// Display the token
	fmt.Println() // Blank line
	color.Green.Println(dbKey)
	fmt.Println() // Blank line

	expiration, dateParseErr := timefmt.Parse(dbKeyExpiration, "%Y-%m-%d %H:%M:%S.%f")
	checkErrCmd(dateParseErr, fmt.Sprintf("%s", dateParseErr), 1)
	if expiration.Year() > 2120 {
		color.Yellow.Println("The above token is permanent (does not expire unless manually revoked).")
	} else {
		color.White.Print("expires in date: ")
		color.Yellow.Println(timefmt.Format(expiration, "%d/%m/%Y (%A %d %B %Y) %H:%M:%S %p"))
	}
	fmt.Println() // Blank line

	os.Exit(0)
}
