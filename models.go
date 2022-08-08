// models.go stores all the table definitions.
// They have to match the Database structure (see sql/*.sql files).
// Hint: Database Null values can be handled using interface{} rather than the actual type.
package main

import (
	"encoding/base64"
	"fmt"
	"reflect"
)

type User struct {
	Id         int         `json:"id"`
	Username   string      `json:"username"`
	Password   string      `json:"password"`
	FirstName  interface{} `json:"first_name"`
	LastName   interface{} `json:"last_name"`
	Email      interface{} `json:"email"`
	Birthday   interface{} `json:"birthday"`
	Picture    interface{} `json:"picture"`
	Phone      interface{} `json:"phone"`
	DateJoined interface{} `json:"date_joined"`
	LastLogin  interface{} `json:"last_login"`
	Role       interface{} `json:"role"`
	IsAdmin    bool        `json:"is_admin"`
	Active     bool        `json:"active"`
}

func (u *User) getFields() []string {
	var fields []string
	t := reflect.TypeOf(u)

	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

func (u *User) getUserById(Id int) error {
	query := "SELECT * FROM User WHERE active=1 AND Id=?;"
	row := db.QueryRow(query, Id)
	err := row.Scan(&u.Id, &u.Username, &u.Password, &u.FirstName, &u.LastName, &u.Email, &u.Birthday, &u.Picture, &u.Phone, &u.DateJoined, &u.LastLogin, &u.Role, &u.IsAdmin, &u.Active)

	if u.Picture != nil {
		// Encode image (BLOB) to base64 string for displaying it in the html template
		u.Picture = base64.StdEncoding.EncodeToString([]byte(u.Picture.([]byte)))
	}

	if err != nil {
		//log.Panic("getUser() error selecting User, err:", err)  // log.Panic would cause a crash and doesn't allow to handle with a proper user logout
		return err
	}
	return nil
}

func (u *User) getUserByUsername() error {
	query := "SELECT id, username, password FROM User WHERE active=1 AND username=?;"
	row := db.QueryRow(query, u.Username)
	err := row.Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		fmt.Println("getUser() error selecting User, err:", err)
		return err
	}
	return nil
}
