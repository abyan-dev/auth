package utils

import (
	"regexp"
	"strings"
	"unicode"

	goaway "github.com/TwiN/go-away"
)

func ValidateEmail(email string) (bool, string) {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegexPattern)
	if !re.MatchString(email) {
		if !strings.Contains(email, "@") {
			return false, "Email must contain an '@' symbol"
		}
		if !strings.Contains(email, ".") {
			return false, "Email must contain a '.' symbol"
		}
		if strings.Contains(email, " ") {
			return false, "Email must not contain spaces"
		}
		return false, "Email format is invalid"
	}
	return true, "Email is valid"
}

func ValidateName(name string) (bool, string) {
	const displayNameRegexPattern = `^[a-zA-Z0-9][a-zA-Z0-9 _-]{1,28}[a-zA-Z0-9]$`
	re := regexp.MustCompile(displayNameRegexPattern)

	if len(name) < 3 || len(name) > 30 {
		return false, "Name must be between 3 and 30 characters long"
	}

	if strings.Contains(name, "  ") {
		return false, "Name must not contain consecutive spaces"
	}

	if !re.MatchString(name) {
		if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(name) {
			return false, "Name must start with an alphanumeric character"
		}
		if !regexp.MustCompile(`[a-zA-Z0-9]$`).MatchString(name) {
			return false, "Name must end with an alphanumeric character"
		}
		return false, "Name contains invalid characters"
	}

	if goaway.IsProfane(name) {
		return false, "Name contains profane words"
	}

	return true, "Name is valid"
}

func ValidatePassword(password string) (bool, string) {
	var (
		hasMinLen    = false
		hasUpperCase = false
		hasLowerCase = false
		hasNumber    = false
		hasSpecial   = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasMinLen {
		return false, "Password must be at least 8 characters long"
	}
	if !hasUpperCase {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLowerCase {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, "Password is valid"
}
