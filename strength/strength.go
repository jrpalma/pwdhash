package strength

import (
	"strings"
	"unicode"
)

// PasswordStrength Represents the rules used to check password strength.
type PasswordStrength struct {
	// MinLowerCase The minimum number of lower case characters in the password.
	MinLowerCase uint
	// MinUpperCase The maximum number of upper case characters in the password.
	MinUpperCase uint
	// MinSpecial The minimum number of special characters in the password.
	// Special characters include: ~!@#$%^&*()_+-=
	MinSpecial uint
	// MinNumber The minimum number of digit characters in the password.
	MinDigits uint
	// MinLength The minimum password length.
	MinLength uint
	// MaxLength The maximum password length. This is used to avoid extremely large passwords.
	MaxLength uint
}

// Check Returns true if the password strength pass the strength rules.
func Check(rules PasswordStrength, password string) bool {
	var numLowerCase, numUpperCase, numDigits, numSpecial, length uint

	codePoints := []rune(password)
	for _, codePoint := range codePoints {
		length++
		if unicode.IsLower(codePoint) {
			numLowerCase++
		}
		if unicode.IsUpper(codePoint) {
			numUpperCase++
		}
		if unicode.IsDigit(codePoint) {
			numDigits++
		}
		if isSpecial(codePoint) {
			numSpecial++
		}

	}

	if numLowerCase < rules.MinLowerCase {
		return false
	}
	if numUpperCase < rules.MinUpperCase {
		return false
	}
	if numDigits < rules.MinDigits {
		return false
	}
	if numSpecial < rules.MinSpecial {
		return false
	}
	if length < rules.MinLength {
		return false
	}
	if length > rules.MaxLength {
		return false
	}

	return true
}

const (
	specialChars = "~!@#$%^&*()_+-="
)

func isSpecial(codePoint rune) bool {
	return strings.Contains(specialChars, string(codePoint))
}
