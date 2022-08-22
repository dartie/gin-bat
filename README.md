<img align="right" width="159px" src="https://raw.githubusercontent.com/dartie/gin-bat/master/logo.svg" width="360">

<!--
<p align="center">
    <a href="https://go.dev/"><img src="https://pkg.go.dev/badge/github.com/dartie/gin-bat" alt="PkgGoDev"></a>
</p>
-->

# Gin-bat

Gin-bat adds the missing batteries to [Gin Web Framework](https://github.com/gin-gonic/gin).
[Gin Web Framework](https://github.com/gin-gonic/gin) is the fastest go web framework available, however it doesn't provide the most required features out-of-box. It's maily inspired by Django.

* MVC structure
* Database connection
* Basic routes
* Authentication system
* Project command line management
* Admin user interface

## Command line
Provides the following sub-commands:

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
...
