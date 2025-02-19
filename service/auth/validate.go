package auth

import "regexp"

func ValidateEmail(email string) bool {
	// 1. validate basic email format
	regex := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

	if !regex.MatchString(email) { // nolint: gosimple
		return false
	}

	// 2. validate specific email model
	// TODO
	return true
}
