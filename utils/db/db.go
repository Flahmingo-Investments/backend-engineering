package db

import (
	"fmt"
	"strings"
)

type User struct {
	Name, PhoneNumber, LastOTP string
	LogInStatus                bool
	// To do : timestamp for login imposing auto logout
}

var mapOfUsers map[string]User = make(map[string]User)

func UpdateLoginOTP(phoneNumber, randomOtpCode string) bool {
	user := GetUserDetails(phoneNumber)
	user.LastOTP = randomOtpCode
	mapOfUsers[phoneNumber] = user
	return true
}

func Logout(phoneNumber string) bool {
	user := GetUserDetails(phoneNumber)
	user.LogInStatus = false
	mapOfUsers[phoneNumber] = user
	return true
}

func newUser(name, phoneNumber string) User {
	return User{name, phoneNumber, "", false}
}

func RegisterUser(name, phoneNumber string) bool {
	if GetUserDetails(phoneNumber).PhoneNumber != "" {
		panic("Phone already registered, use another !")
	}

	mapOfUsers[phoneNumber] = newUser(name, phoneNumber)
	return true
}

func GetUserDetails(phoneNumber string) User {
	return mapOfUsers[phoneNumber]
}

func LoginUsingOTP(phoneNumber, otpCode string) bool {
	user := GetUserDetails(phoneNumber)
	if strings.Compare(user.LastOTP, otpCode) == 0 {
		user.LogInStatus = true
		mapOfUsers[phoneNumber] = user
		return true
	}
	return false
}

func PrintAll() {
	fmt.Println(mapOfUsers)
}

