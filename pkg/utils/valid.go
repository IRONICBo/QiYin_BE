package utils

import "regexp"

// IsValidEmail check if the email is valid.
func IsValidEmail(email string) bool {
	// Regular expression for basic email validation
	// This regex pattern is a simplified version and may not cover all edge cases.
	// You can use more comprehensive patterns depending on your specific needs.
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	return regexp.MustCompile(emailRegex).MatchString(email)
}
