package models

import "errors"

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

var users = []User{
	{ID: "1", Username: "admin", Password: "adminpassword", Role: "admin", Department: "IT"},
	{ID: "2", Username: "hrmanager", Password: "hrpassword", Role: "manager", Department: "HR"},
	{ID: "3", Username: "itstaff", Password: "itpassword", Role: "staff", Department: "IT"},
	{ID: "4", Username: "employee1", Password: "emppassword1", Role: "employee", Department: "Sales"},
	{ID: "5", Username: "employee2", Password: "emppassword2", Role: "employee", Department: "Marketing"},
}

func GetUserByUsername(username string) (User, error) {
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}
