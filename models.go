package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
)

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Birthday   string `json:"birthday"`
	Picture    string `json:"picture"`
	Phone      string `json:"phone"`
	DateJoined string `json:"date_joined"`
	LastLogin  string `json:"last_login"`
	Role       string `json:"role"`
	IsAdmin    bool   `json:"is_admin"`
	Active     bool   `json:"active"`
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
	query := "SELECT * FROM User WHERE Id=?;"
	row := db.QueryRow(query, Id)
	err := row.Scan(&u.Id, &u.Username, &u.Password, &u.FirstName, &u.LastName, &u.Email, &u.Birthday, &u.Picture, &u.Phone, &u.DateJoined, &u.LastLogin, &u.Role, &u.IsAdmin, &u.Active)

	// Encode image (BLOB) to base64 string for displaying it in the html template
	u.Picture = base64.StdEncoding.EncodeToString([]byte(u.Picture))

	if err != nil {
		log.Panic("getUser() error selecting User, err:", err)
		return err
	}
	return nil
}

func (u *User) getUserByUsername() error {
	query := "SELECT id, username, password FROM User WHERE username=?;"
	row := db.QueryRow(query, u.Username)
	err := row.Scan(&u.Id, &u.Username, &u.Password)
	if err != nil {
		fmt.Println("getUser() error selecting User, err:", err)
		return err
	}
	return nil
}
