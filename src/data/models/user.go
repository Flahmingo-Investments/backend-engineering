package models

type User struct {
	Base
	Name     string
	Phone    string
	Verified bool
	Logged   bool
}
