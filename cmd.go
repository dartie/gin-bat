package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gookit/color"
	"github.com/hellflame/argparse"
	"github.com/howeyc/gopass"
	"github.com/itchyny/timefmt-go"
	"github.com/scylladb/termtables"
	"github.com/soypat/rebed"
	str2duration "github.com/xhit/go-str2duration/v2"
	"golang.org/x/crypto/bcrypt"
)

var appName = "MyWebApp"
var Version = "1.0"
var cwd, errCwd = os.Getwd()

/* Arguments command line */
var parser = argparse.NewParser(appName+" "+Version, ``, nil)
var createProjectCommand = parser.AddCommand("create-project", "Create a new project", nil)
var createUserCommand = parser.AddCommand("create-user", "Create a new user", nil)
var updateUserCommand = parser.AddCommand("update-user", "Update an existing user (Passwords can be changed with change-password command)", nil)
var deleteUserCommand = parser.AddCommand("delete-user", "Delete an existing user", nil)
var changePasswordCommand = parser.AddCommand("change-password", "Change password for an existing user", nil)
var createUserTokenCommand = parser.AddCommand("create-token", "Create a new user", nil)
var displayUserTokenCommand = parser.AddCommand("display-token", "Display token for a given user", nil)
var generateSqlFilesCommand = parser.AddCommand("generate-sql", "Generate the sql files for all compatible databases", &argparse.ParserConfig{DisableDefaultShowHelp: true})
var runserverCommand = parser.AddCommand("runserver", "Run the server", &argparse.ParserConfig{DisableDefaultShowHelp: true})

// Variables for create-project command
var createProjectPath *string

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
var updateUserFirstname *string
var updateUserLastname *string
var updateUserEmail *string
var updateUserPhone *string
var updateUserBirthday *string
var updateUserPicture *string
var updateUserRole *string
var updateUserAdmin *bool
var updateUserActive *bool

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

// Variables for display-token command
var generateSqlFilesDbType *string

// Variables for runserver command
var host *string
var port *string

func initcmd() {
	// Arguments for createUserCommand
	createProjectPath = createProjectCommand.String("p", "project-name", &argparse.Option{Positional: true})

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
	updateUserFirstname = updateUserCommand.String("n", "first-name", nil)
	updateUserLastname = updateUserCommand.String("s", "last-name", nil)
	updateUserEmail = updateUserCommand.String("e", "email", nil)
	updateUserPhone = updateUserCommand.String("t", "telephone", nil)
	updateUserBirthday = updateUserCommand.String("b", "birthday", &argparse.Option{Help: "Format YYYY-MM-DD"})
	updateUserPicture = updateUserCommand.String("pic", "picture", nil)
	updateUserRole = updateUserCommand.String("r", "role", nil)
	updateUserAdmin = updateUserCommand.Flag("A", "admin", nil)
	updateUserActive = updateUserCommand.Flag("a", "active", nil)

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

	if createProjectCommand.Invoked {
		createProject()
	}

	if createUserCommand.Invoked {
		createUser()
	}

	if updateUserCommand.Invoked {
		updateUser()
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

	if generateSqlFilesCommand.Invoked {
		generateSqlFiles()
	}

	if displayUserTokenCommand.Invoked {
		displayUserToken()
	}
}

func createProject() {
	projectDir := filepath.Dir(*createProjectPath)
	projectName := filepath.Base(*createProjectPath)
	if _, err := os.Stat(*createProjectPath); errors.Is(err, os.ErrNotExist) {
		errMkdir := os.MkdirAll(*createProjectPath, os.ModePerm)
		checkErrCmd(errMkdir, fmt.Sprintf("%s", errMkdir), 1)
	}

	errRebed := rebed.Write(fsProjectFiles, *createProjectPath)
	checkErrCmd(errRebed, fmt.Sprintf("%s", errRebed), 1)

	/* Set name for go project */
	gomodFilePath := filepath.Join(*createProjectPath, "go.mod")
	input, errReadFile := os.ReadFile(gomodFilePath)
	checkErrCmd(errReadFile, fmt.Sprintf("%s", errReadFile), 1)

	output := bytes.Replace(input, []byte("module ginbat"), []byte("module "+projectName), -1)

	errWriteFile := os.WriteFile(gomodFilePath, output, 0666)
	checkErrCmd(errWriteFile, fmt.Sprintf("%s", errWriteFile), 1)

	if projectDir == "." {
		projectDir, _ = os.Getwd()
	}

	fmt.Println()
	color.Green.Printf("Project \"%s\" created succefully in \"%s\"\n", projectName, projectDir)
	fmt.Println()

	os.Exit(0)
}

func createUser() {
	sqlInsertString := `INSERT INTO "Users" ( 
		username, password, first_name, last_name, email, birthday, picture, phone, date_joined, last_login, role, is_admin, active
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

	// Check inputs
	if *createUserUser == "" {
		color.Cyan.Print("Enter username:\n>")
		*createUserUser = readStdin()
	}

	// Check if new profile already exists
	tokenKey := strings.ToUpper(*createUserUser)
	query := `SELECT username FROM Users where UPPER(username) = ?`
	row := db.QueryRow(query, tokenKey)
	var dbid interface{}
	queryErr := row.Scan(&dbid)
	if queryErr == nil {
		log.Fatalf(color.Red.Sprintf("User \"%s\" already exists. Exiting.", *createUserUser))
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
	checkErrCmd(hashedPasswordErr, fmt.Sprintf("%s", hashedPasswordErr), 1)

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
	var readPicErr error
	if *createUserPicture != "" {
		//Validate path
		if *createUserPicture != "" {
			if _, err := os.Stat(*createUserPicture); errors.Is(err, os.ErrNotExist) {
				log.Fatal(color.Red.Sprintf("File \"%s\" doesn't exists", *createUserPicture))
			}
		}
		userPictureFile, fileErr := os.Open(*createUserPicture)
		checkErrCmd(fileErr, fmt.Sprintf("%s", fileErr), 1)
		defer userPictureFile.Close()
		profileData, readPicErr = io.ReadAll(userPictureFile)
		checkErrCmd(readPicErr, fmt.Sprintf("%s", readPicErr), 1)
	} else {
		profileData = nil
	}

	if *createUserPhone == "" {
		color.Cyan.Print("Enter phone:\n>")
		*createUserPhone = readStdin()
	}

	var adminInputString = "NO"
	var createUserAdminInt = 0 // for query insert
	if !*createUserAdmin {
		var adminInput string
		color.Cyan.Print("Is the user Admin? (Y\\n):\n>")
		adminInput = readStdin()

		if strings.ToUpper(adminInput) == "Y" {
			*createUserAdmin = true
			adminInputString = "YES"
			createUserAdminInt = 1
		}
	} else {
		adminInputString = "YES"
		createUserAdminInt = 1
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
	checkErrCmd(err, fmt.Sprintf("%s", err), 1)
	sqlResult, sqlErr := sqlCommand.Exec(*createUserUser, hashedPassword, *createUserFirstname, *createUserLastname, *createUserEmail, *createUserBirthday, profileData, *createUserPhone, nowSqliteFormat(), "", "", createUserAdminInt, 1)
	var status int64
	var message string
	if sqlErr == nil {
		recordId, err := sqlResult.LastInsertId() // doesn't work with Postgres.
		if err != nil {
			recordId = 0
		}
		message = fmt.Sprintf("User \"%s\" has been created successfully", *createUserUser)
		status = recordId
		color.Green.Println(message)

	} else {
		message = fmt.Sprintf("Issues creating User \"%s\" : %s", *createUserUser, sqlErr)
		status = 0
		color.Red.Println(message)

	}
	os.Exit(int(status))
}

func updateUser() {
	// Check inputs
	if *updateUserUser == "" {
		color.Cyan.Print("Enter username:\n>")
		*updateUserUser = readStdin()
	}

	// Check if the profile exists
	tokenKey := strings.ToUpper(*updateUserUser)
	query := `SELECT Id, username, first_name, last_name, email, birthday, phone, picture, date_joined, last_login, role, is_admin, active FROM "Users" WHERE UPPER(username) = $1`
	row := db.QueryRow(query, tokenKey)
	var dbid int
	var dbisadmin, dbactive bool
	var dbusername, dbfirstname, dblastname, dbemail, dbbirthday, dbphone, dbdatejoined, dblastlogin, dbrole string
	var dbpicture interface{}
	queryErr := row.Scan(&dbid, &dbusername, &dbfirstname, &dblastname, &dbemail, &dbbirthday, &dbphone, &dbpicture, &dbdatejoined, &dblastlogin, &dbrole, &dbisadmin, &dbactive)

	if queryErr != nil {
		log.Fatalf(color.Red.Sprintf("User \"%s\" Doesn't exist. Exiting.", *updateUserUser))
	}

	color.Yellow.Println("Type the new value or press Enter for confirming the stored value.")
	color.Yellow.Println("Enter a space for providing a null value.")
	fmt.Println()

	if *updateUserFirstname == "" {
		color.Cyan.Print("Enter new First Name:\n")
		color.Yellow.Printf("Current value: \"%s\"", dbfirstname)
		color.Cyan.Print("\n>")
		InputUpdateUserFirstname := strings.ReplaceAll(readStdinOriginal(), "\n", "")

		if InputUpdateUserFirstname == "" {
			// leave the previous value
			*updateUserFirstname = dbfirstname
		} else if strings.HasPrefix(InputUpdateUserFirstname, " ") && strings.TrimSpace(InputUpdateUserFirstname) == "" {
			// overwrite previous value with null
			*updateUserFirstname = ""
		} else {
			*updateUserFirstname = InputUpdateUserFirstname
		}
	}

	if *updateUserLastname == "" {
		color.Cyan.Print("Enter new Last Name:\n")
		color.Yellow.Printf("Current value: \"%s\"", dblastname)
		color.Cyan.Print("\n>")
		InputUpdateUserLastname := strings.ReplaceAll(readStdinOriginal(), "\n", "")

		if InputUpdateUserLastname == "" {
			// leave the previous value
			*updateUserLastname = dblastname
		} else if strings.HasPrefix(InputUpdateUserLastname, " ") && strings.TrimSpace(InputUpdateUserLastname) == "" {
			// overwrite previous value with null
			*updateUserLastname = ""
		} else {
			*updateUserLastname = InputUpdateUserLastname
		}
	}

	if *updateUserEmail == "" {
		color.Cyan.Print("Enter new email:\n")
		color.Yellow.Printf("Current value: \"%s\"", dbemail)
		color.Cyan.Print("\n>")
		InputUpdateUserEmail := strings.ReplaceAll(readStdinOriginal(), "\n", "")

		if InputUpdateUserEmail == "" {
			// leave the previous value
			*updateUserEmail = dbemail
		} else if strings.HasPrefix(InputUpdateUserEmail, " ") && strings.TrimSpace(InputUpdateUserEmail) == "" {
			// overwrite previous value with null
			*updateUserEmail = ""
		} else {
			*updateUserEmail = InputUpdateUserEmail
		}
	}

	if *updateUserBirthday == "" {
		color.Cyan.Print("Enter new birthday (YYYY-MM-DD):\n")
		color.Yellow.Printf("Current value: \"%s\"", dbbirthday)
		color.Cyan.Print("\n>")
		InputUpdateUserBirthday := strings.ReplaceAll(readStdinOriginal(), "\n", "")

		if InputUpdateUserBirthday == "" {
			// leave the previous value
			*updateUserBirthday = dbbirthday
		} else if strings.HasPrefix(InputUpdateUserBirthday, " ") && strings.TrimSpace(InputUpdateUserBirthday) == "" {
			// overwrite previous value with null
			*updateUserBirthday = ""
		} else {
			*updateUserBirthday = InputUpdateUserBirthday
		}
	}

	if *updateUserBirthday != "" {
		// Validate date input
		//re := regexp.MustCompile(`(0?[1-9]|[12][0-9]|3[01])[-/](0?[1-9]|1[012])[-/]((19|20)\d\d)`)
		re := regexp.MustCompile(`((19|20)\d\d)[-/](0?[1-9]|1[012])[-/](0?[1-9]|[12][0-9]|3[01])`)
		if !re.MatchString(*updateUserBirthday) {
			log.Fatalf(color.Red.Sprintf("Invalid date format: %s", *updateUserBirthday))
		}

		dateObject, err := dateparse.ParseAny(*updateUserBirthday)
		if err != nil {
			log.Fatalf(color.Red.Sprintf("Invalid date: %s", *updateUserBirthday))
		}

		*updateUserBirthday = timefmt.Format(dateObject, "%Y-%m-%d")
	}

	var profileData interface{} //[]byte
	var updateUserPictureChange string
	if *updateUserPicture == "" {
		var dbpicturePreview string
		if dbpicture == "" || dbpicture == nil {
			dbpicturePreview = "<no image>"
		} else {
			dbpicturePreview = "Image stored"
		}
		color.Cyan.Printf("Enter new picture path:\n")
		color.Yellow.Printf("Current value: \"%s...\"", dbpicturePreview)
		color.Cyan.Print("\n>")
		InputUpdateUserPicture := strings.ReplaceAll(readStdinOriginal(), "\n", "")

		if InputUpdateUserPicture == "" {
			// leave the previous value
			profileData = dbpicture
		} else if strings.HasPrefix(InputUpdateUserPicture, " ") && strings.TrimSpace(InputUpdateUserPicture) == "" {
			// overwrite previous value with null
			profileData = nil
			if dbpicture != nil {
				updateUserPictureChange = "Picture removed"
			}
		} else {
			updateUserPictureChange = "New picture uploaded"
			*updateUserPicture = InputUpdateUserPicture

			//Validate path
			if _, err := os.Stat(*updateUserPicture); errors.Is(err, os.ErrNotExist) {
				log.Fatal(color.Red.Sprintf("File \"%s\" doesn't exists", *updateUserPicture))
			}
			userPictureFile, fileErr := os.Open(*updateUserPicture)
			checkErrCmd(fileErr, fmt.Sprintf("%s", fileErr), 1)
			defer userPictureFile.Close()
			var readPicErr error
			profileData, readPicErr = io.ReadAll(userPictureFile)
			checkErrCmd(readPicErr, fmt.Sprintf("%s", readPicErr), 1)
		}
	}

	if *updateUserPhone == "" {
		color.Cyan.Print("Enter phone:\n>")
		color.Yellow.Printf("Current value: \"%s\"", dbphone)
		color.Cyan.Print("\n>")
		InputUpdateUserPhone := strings.ReplaceAll(readStdinOriginal(), "\n", "")

		if InputUpdateUserPhone == "" {
			// leave the previous value
			*updateUserPhone = dbphone
		} else if strings.HasPrefix(InputUpdateUserPhone, " ") && strings.TrimSpace(InputUpdateUserPhone) == "" {
			// overwrite previous value with null
			*updateUserPhone = ""
		} else {
			*updateUserPhone = InputUpdateUserPhone
		}
	}

	var updateUserAdminInt int
	if !*updateUserAdmin {
		color.Cyan.Print("Is the user Admin? (Y\\n):\n>")
		color.Yellow.Printf("Current value: \"%v\"", dbisadmin)
		color.Cyan.Print("\n>")
		InputUpdateUserAdmin := readStdin()

		if InputUpdateUserAdmin == "" {
			// leave the previous value
			*updateUserAdmin = dbisadmin
		} else if strings.ToUpper(InputUpdateUserAdmin) == "Y" || strings.ToUpper(InputUpdateUserAdmin) == "YES" {
			// overwrite previous value with null
			*updateUserAdmin = true
		} else {
			*updateUserAdmin = false
		}
	}

	if *updateUserAdmin {
		updateUserAdminInt = 1
	} else {
		updateUserAdminInt = 0
	}

	var updateUserActiveInt int
	if !*updateUserActive {
		color.Cyan.Print("Is the user Active? (Y\\n):\n>")
		color.Yellow.Printf("Current value: \"%v\"", dbactive)
		color.Cyan.Print("\n>")
		InputUpdateUserActive := readStdin()

		if InputUpdateUserActive == "" {
			// leave the previous value
			*updateUserActive = dbactive
		} else if strings.ToUpper(InputUpdateUserActive) == "Y" || strings.ToUpper(InputUpdateUserActive) == "YES" {
			// overwrite previous value with null
			*updateUserActive = true
		} else {
			*updateUserActive = false
		}
	}

	if *updateUserActive {
		updateUserActiveInt = 1
	} else {
		updateUserActiveInt = 0
	}

	// Ask for username change
	color.Cyan.Print("Enter new Username:\n")
	color.Yellow.Printf("Current value: \"%s\"", dbusername)
	color.Cyan.Print("\n>")
	InputUpdateUserUser := strings.ReplaceAll(readStdinOriginal(), "\n", "")

	if InputUpdateUserUser == "" {
		// leave the previous value
		*updateUserUser = dbusername
	} else if strings.HasPrefix(InputUpdateUserUser, " ") && strings.TrimSpace(InputUpdateUserUser) == "" {
		// overwrite previous value with null
		*updateUserUser = ""
	} else {
		*updateUserUser = InputUpdateUserUser
	}

	/* Print Summary */
	// Get the changes
	var userWarning string
	var updateUserUsernameChange string
	if *updateUserUser != dbusername {
		updateUserUsernameChange = fmt.Sprintf("%s -> %s", dbusername, *updateUserUser)
		userWarning = fmt.Sprintf("Please note, user %s is no longer available and you will need to use %s for referring to this user", *updateUserUser, dbusername)
	}
	var updateUserFirstnameChange string
	if *updateUserFirstname != dbfirstname {
		updateUserFirstnameChange = fmt.Sprintf("%s -> %s", dbfirstname, *updateUserFirstname)
	}
	var updateUserLastnameChange string
	if *updateUserFirstname != dbfirstname {
		updateUserLastnameChange = fmt.Sprintf("%s -> %s", dblastname, *updateUserLastname)
	}
	var updateUserEmailChange string
	if *updateUserEmail != dbemail {
		updateUserEmailChange = fmt.Sprintf("%s -> %s", dbemail, *updateUserEmail)
	}
	var updateUserBirthdayChange string
	if *updateUserBirthday != dbbirthday {
		updateUserBirthdayChange = fmt.Sprintf("%s -> %s", dbbirthday, *updateUserBirthday)
	}
	var updateUserPhoneChange string
	if *updateUserPhone != dbphone {
		updateUserPhoneChange = fmt.Sprintf("%s -> %s", dbphone, *updateUserPhone)
	}
	var updateUserAdminChange string
	if *updateUserAdmin != dbisadmin {
		updateUserAdminChange = fmt.Sprintf("%v -> %v", dbisadmin, *updateUserAdmin)
	}
	var updateUserActiveChange string
	if *updateUserActive != dbisadmin {
		updateUserActiveChange = fmt.Sprintf("%v -> %v", dbactive, *updateUserActive)
	}

	table := termtables.CreateTable()
	table.AddHeaders("Info", "DB Field", "Value", "Changes")
	table.AddRow("Username", "username", *updateUserUser, updateUserUsernameChange)
	table.AddRow("First Name", "first_name", *updateUserFirstname, updateUserFirstnameChange)
	table.AddRow("Last Name", "last_name", *updateUserLastname, updateUserLastnameChange)
	table.AddRow("Email", "email", *updateUserEmail, updateUserEmailChange)
	table.AddRow("Birthday", "birthday", *updateUserBirthday, updateUserBirthdayChange)
	table.AddRow("Profile Picture", "picture", *updateUserPicture, updateUserPictureChange)
	table.AddRow("Phone number", "phone", *updateUserPhone, updateUserPhoneChange)
	table.AddRow("Administrator", "is_admin", *updateUserAdmin, updateUserAdminChange)
	table.AddRow("Active", "active", *updateUserActive, updateUserActiveChange)
	color.Cyan.Println(table.Render())

	var confirmSummary string
	color.Yellow.Println("Update the user with the above info? [Y\\n]\n>")
	confirmSummary = readStdin()

	if strings.ToUpper(confirmSummary) != "Y" {
		color.Red.Println("User aborted")
		os.Exit(0)
	}

	// DB Update
	sqlUpdateString := `UPDATE "Users" 
	SET username=$1, first_name=$2, last_name=$3, email=$4, birthday=$5, picture=$6, phone=$7, role=$8, is_admin=$9, active=$10
	WHERE id=$11;
	`
	sqlCommand, err := db.Prepare(sqlUpdateString)
	checkErrCmd(err, fmt.Sprintf("%s", err), 1)
	_, sqlErr := sqlCommand.Exec(*updateUserUser, *updateUserFirstname, *updateUserLastname, *updateUserEmail, *updateUserBirthday, profileData, *updateUserPhone, "", updateUserAdminInt, updateUserActiveInt, dbid)
	var status int
	var message string
	if sqlErr == nil {
		message = fmt.Sprintf("User \"%s\" - Id = \"%d\" has been update successfully", *updateUserUser, dbid)
		status = dbid
		color.Green.Println(message)

	} else {
		message = fmt.Sprintf("Issues updating User \"%s\" : %s", *createUserUser, sqlErr)
		status = 0
		color.Red.Println(message)
	}

	color.Magenta.Println(userWarning)

	os.Exit(int(status))
}

func deleteUser() {
	sqlDeleteString := `DELETE FROM Users WHERE username=?`

	if *deleteUserUser == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the user delete (username)\n>")
		*deleteUserUser = readStdinOriginal()
	}

	// Check if the user exists
	query := `SELECT Id, username, first_name, last_name, email, birthday, phone, date_joined, last_login, role, is_admin, active FROM "Users" WHERE username=$1;`
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
		log.Fatalf(color.Red.Sprintf("User \"%s\" does not exist.", *deleteUserUser))
	}

	// Ask confirm
	var confirm string
	color.Yellow.Printf("Are you sure you want to delete \"%s\" user? [Y\\n]\n>", *deleteUserUser)
	confirm = readStdin()

	if strings.ToUpper(confirm) == "Y" {
		sqlCommand, err := db.Prepare(sqlDeleteString)
		checkErrCmd(err, fmt.Sprintf("%s", err), 1)
		_, sqlErr := sqlCommand.Exec(*deleteUserUser)
		checkErrCmd(sqlErr, fmt.Sprintf("%s", sqlErr), 1)
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
	query := `SELECT Id, username, first_name, last_name, email, birthday, phone, date_joined, last_login, role, is_admin, active FROM "Users" WHERE username=$1`
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
	checkErrCmd(hashedPasswordErr, fmt.Sprintf("%s", hashedPasswordErr), 1)

	// Update User with new password
	updatePasswordSql := `UPDATE "Users" SET password=$1 WHERE id=$2`

	sqlCommand, err := db.Prepare(updatePasswordSql)
	checkErrCmd(err, fmt.Sprintf("%s", err), 1)
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
	query := `SELECT Id, username, first_name, last_name, email, birthday, phone, date_joined, last_login, role, is_admin, active FROM "Users" WHERE username=$1;`
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

	query = `SELECT user_id FROM "AuthToken" WHERE user_id=$1`
	row = db.QueryRow(query, u.Id)
	selectErr = row.Scan(&dbId)
	if selectErr == nil {
		sqlInsertString = `UPDATE "AuthToken" SET token_key=$1, created=$2, expiration=$3 WHERE user_id=$4`
	} else {
		// A record is already present -> UPDATE
		sqlInsertString = `INSERT INTO "AuthToken" ( 
			token_key, created, expiration, user_id
			)
			VALUES
			($1, $2, $3, $4)
			`
	}

	// Execute Query
	sqlCommand, err := db.Prepare(sqlInsertString)
	checkErrCmd(err, fmt.Sprintf("%s", err), 1)
	_, sqlErr := sqlCommand.Exec(tokenString, nowSqliteFormat(), expirationDBFormat, u.Id)
	checkErrCmd(sqlErr, fmt.Sprintf("%s", sqlErr), 1)

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

func generateSqlFiles() {
	sqlFolder := "sql"
	pgExt := "-pg.sql"
	fileOpen, err := os.Open(sqlFolder)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOpen.Close()
	files, _ := fileOpen.Readdir(0)

	for _, f := range files {
		if strings.HasSuffix(f.Name(), pgExt) {
			continue
		}
		ext := filepath.Ext(f.Name())
		filename := f.Name()[0 : len(f.Name())-len(ext)]
		InputFileFullpath := filepath.Join(cwd, sqlFolder, f.Name())
		OutputFileFullpath := filepath.Join(cwd, sqlFolder, filename+pgExt)
		content, err := os.ReadFile(InputFileFullpath)
		checkErrCmd(err, fmt.Sprint(err), 1)

		// Postgres
		newContent := regexp.MustCompile(`(?i)BLOB`).ReplaceAllString(string(content), "bytea")
		newContent = regexp.MustCompile(`(?i)DATETIME`).ReplaceAllString(string(newContent), "DATE")

		err = os.WriteFile(OutputFileFullpath, []byte(newContent), 0644)
		checkErrCmd(err, fmt.Sprint(err), 1)
	}

	// Print ok
	fmt.Println() // Blank line
	color.Green.Println("Files for Postgres generated correctly")
	fmt.Println() // Blank line
	os.Exit(0)
}

func displayUserToken() {
	if *displayUserTokenUser == "" {
		// Ask username to remove
		color.Cyan.Printf("Type the user you want to display the token\n>")
		*displayUserTokenUser = readStdin()
	}

	query := `SELECT token_key, expiration
	FROM "AuthToken"
	INNER JOIN "Users"
	ON "AuthToken".user_id = "Users".Id
	WHERE "Users".username = $1;
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
