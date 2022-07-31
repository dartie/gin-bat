package main

import "fmt"

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
