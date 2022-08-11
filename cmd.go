package main

import (
	"fmt"

	"github.com/hellflame/argparse"
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

func initcmd() {
	// Arguments for createUserCommand
	createUserCommand.String("u", "user", nil)
	createUserCommand.String("p", "password", nil)
	createUserCommand.String("n", "first-name", nil)
	createUserCommand.String("s", "last-name", nil)
	createUserCommand.String("e", "email", nil)
	createUserCommand.String("t", "telephone", nil)
	createUserCommand.String("b", "birthday", nil)
	createUserCommand.String("pic", "picture", nil)
	createUserCommand.String("r", "role", nil)
	createUserCommand.Flag("a", "admin", nil)

	// Arguments for updateUserCommand
	updateUserCommand.String("u", "user", nil)
	updateUserCommand.String("p", "password", nil)
	updateUserCommand.String("n", "first-name", nil)
	updateUserCommand.String("s", "last-name", nil)
	updateUserCommand.String("e", "email", nil)
	updateUserCommand.String("t", "telephone", nil)
	updateUserCommand.String("b", "birthday", nil)
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
}
