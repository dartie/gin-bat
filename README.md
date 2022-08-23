<img align="right" width="159px" src="https://raw.githubusercontent.com/dartie/gin-bat/master/logo.svg" width="360">

<!--
<p align="center">
    <a href="https://go.dev/"><img src="https://pkg.go.dev/badge/github.com/dartie/gin-bat" alt="PkgGoDev"></a>
</p>
-->

# Gin-bat

Gin-bat adds the missing batteries to [Gin Web Framework](https://github.com/gin-gonic/gin).
[Gin Web Framework](https://github.com/gin-gonic/gin) is the fastest go web framework available, however it doesn't provide the most required features out-of-box. It's maily inspired by Django.

Features:
* MVC structure
* Database connection
* Basic routes
* Authentication system
* Project command line management
* Admin user interface

## Contents
- [Gin-bat](#gin-bat)
  - [Contents](#contents)
  - [Command line](#command-line)
  - [Project structure](#project-structure)
  - [Getting started](#getting-started)
  - [Customization](#customization)
    - [Logo](#logo)
    - [Title](#title)
    - [Imports](#imports)
    - [Create a new page](#create-a-new-page)
    - [Raise a html alert](#raise-a-html-alert)
    - [Redirect](#redirect)

## Command line
Provides the following sub-commands:

* `create-project` : create a new project, copies all the required files for building a new server, to the given path
* `create-user` : create a new user profile
* `update-user` : modify an existing profile
* `delete-user` : delete a user
* `change-password` : change password access for a user
* `create-token` : create a token for an existing user, to enable api usage
* `display-token` : display the token for a given user
* `runserver` : run the server

## Project structure
Once the project has been created, the following files are available for working on your web-application:

* `urls.go` : defines the route
* `controller.go` : defines the back-end source code (handlers)
* `models.go` : defines database tables golang objects (`User` struct with its methods are provided by default)
* `common.go` : includes function which are shared in the project. You add more there.

The remaining go files are supposed to be untouched, unless there are behaviours you wish to change to the core:
* `main.go`
* `cmd.go`
* `built-in.go`
* `authentication.go`


## Getting started
1. Create a project with
    ```bash
    ginbat create-project <project-path>
    ```
    
    `<project-path>` will contain all the necessary files for running the server

1. Run the server
    ```bash
    go run . runserver 
    ```

    Alternatively provide `-H` and/or `-P` for specifying respectively Host and Port

    ```bash
    go run . runserver -H 0.0.0.0 -P 8080
    ```


## Customization

### Logo
Modify  `templates/structure/topbar.html`

### Title
Modify  `templates/structure/header.html`

### Imports
For convenience, all imports are located in `templates/structure/imports.html` - edit it for adding more.

### Create a new page
1. Create html file in `templates` folder
1. For giving the same style, include the "structure" html files using the `template` keyword (remember to include the dot `.` character for passing the variables from the first html template to the nested ones).
    ```html
    {{ template "Header"  . }}
    {{ template "Imports" . }}
    {{ template "Middle"  . }}
    {{ template "Topbar"  . }}
    <!-- Content -->
    <h1>Home</h1>
    {{ template "Footer"  . }}
    </body>
    </html>
    ```
1. Customize the content below the comment `<!-- Content -->` leaving the inclusions untouched.

### Raise a html alert
ginbat ships a javascript function `raiseAlert()` (defined in `templates/structure/alert.html`) which renders [bootstrap5 alerts](https://getbootstrap.com/docs/5.2/components/alerts/)

From the controller, return
```go
func myHandler(c *gin.Context) {
  message := "My message"
  status := "0"  // successful status

  c.HTML(http.StatusOK, "page.html", gin.H{"Feedback": map[string]string{message: status}})
}
```

Accepted status are:
* `"0"` : :green_circle: success :arrow_right: Green
* `"1"` : :red_circle: danger :arrow_right: Red
* `"2"` : :orange_circle: warning :arrow_right: Orange

> Hint: multiple (unique) messages can be provided in the map.

### Redirect
Gin redirect works fine, however the page rendering doesn't allow to pass parameters. As workaround, you can render the page you wish and pass `"Url"` parameter for adjusting the url in the browser.

This is defined in `templates/structure/footer.html`

From the controller, return
```go
func myHandler(c *gin.Context) {
  c.HTML(http.StatusOK, "page.html", gin.H{"Url": "/"})
}
```
