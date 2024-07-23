package utils

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
		message  string
	}{
		// Valid emails
		{"example@example.com", true, "Email is valid"},
		{"user.name+tag+sorting@example.com", true, "Email is valid"},
		{"user_name@sub.domain.com", true, "Email is valid"},
		{"user-name@domain.co.uk", true, "Email is valid"},

		// Invalid emails
		{"plainaddress", false, "Email must contain an '@' symbol"},
		{"@missingusername.com", false, "Email format is invalid"},
		{"username@.com", false, "Email format is invalid"},
		{"username@com", false, "Email must contain a '.' symbol"},
		{"username@domain.c", false, "Email format is invalid"},
		{"username@domain.com ", false, "Email must not contain spaces"},
		{" username@domain.com", false, "Email must not contain spaces"},
		{"username@ domain.com", false, "Email must not contain spaces"},
		{"user name@domain.com", false, "Email must not contain spaces"},
		{"username@domain.com.", false, "Email format is invalid"},
	}

	for _, test := range tests {
		result, message := ValidateEmail(test.email)
		if result != test.expected || message != test.message {
			t.Errorf("ValidateEmail(%q) = (%v, %q); want (%v, %q)", test.email, result, message, test.expected, test.message)
		}
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		message  string
	}{
		// Valid names
		{"JohnDoe", true, "Name is valid"},
		{"Alice123", true, "Name is valid"},
		{"John_Doe", true, "Name is valid"},
		{"John-Doe", true, "Name is valid"},
		{"John Doe", true, "Name is valid"},
		{"J", false, "Name must be between 3 and 30 characters long"},
		{"Jo", false, "Name must be between 3 and 30 characters long"},
		{"", false, "Name must be between 3 and 30 characters long"},

		// Invalid names
		{"John  Doe", false, "Name must not contain consecutive spaces"},
		{" JohnDoe", false, "Name must start with an alphanumeric character"},
		{"JohnDoe ", false, "Name must end with an alphanumeric character"},
		{"John@Doe", false, "Name contains invalid characters"},
		{"Jo", false, "Name must be between 3 and 30 characters long"},
		{"John Doe with a very long name that exceeds the maximum limit", false, "Name must be between 3 and 30 characters long"},
		{"fuck", false, "Name contains profane words"}, // Replace "badword" with an actual profane word to test
	}

	for _, test := range tests {
		result, message := ValidateName(test.name)
		if result != test.expected || message != test.message {
			t.Errorf("ValidateName(%q) = (%v, %q); want (%v, %q)", test.name, result, message, test.expected, test.message)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
		message  string
	}{
		// Valid passwords
		{"Password1!", true, "Password is valid"},
		{"Valid@123", true, "Password is valid"},
		{"Complex#Pass123", true, "Password is valid"},

		// Invalid passwords
		{"short", false, "Password must be at least 8 characters long"},
		{"NoNumbers!", false, "Password must contain at least one number"},
		{"NOLOWERCASE1!", false, "Password must contain at least one lowercase letter"},
		{"nouppercase1!", false, "Password must contain at least one uppercase letter"},
		{"NoSpecial1", false, "Password must contain at least one special character"},
	}

	for _, test := range tests {
		result, message := ValidatePassword(test.password)
		if result != test.expected || message != test.message {
			t.Errorf("ValidatePassword(%q) = (%v, %q); want (%v, %q)", test.password, result, message, test.expected, test.message)
		}
	}
}
