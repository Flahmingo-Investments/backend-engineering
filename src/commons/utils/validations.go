package commons

import "regexp"

func IsPhoneNumber(phone string) bool {
	match, _ := regexp.MatchString(`^\+\d{11,}`, phone)
	return match
}

func IsEmpty(text string) bool {
	return len(text) == 0
}
