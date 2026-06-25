package utils

// PasswordMinLength is the minimum allowed password length.
const PasswordMinLength = 8

// PasswordMaxLength is bcrypt's practical maximum.
const PasswordMaxLength = 72

func IsPasswordLengthValid(password string) bool {
	length := len(password)

	return length >= PasswordMinLength &&
		length <= PasswordMaxLength
}
