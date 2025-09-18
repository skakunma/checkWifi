package password

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"math/big"
)

type Password = string
type Passwords = []string

const (
	defaultPasswordLength = 8
	numPasswords          = 1000
)

// Character sets for password generation
const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
)

// GenerateRandomPasswords generates 1000 random passwords with specified length
// If length is 0 or negative, uses default length of 12 characters
func GenerateRandomPasswords(length int) (Passwords, error) {
	if length < 0 {
		return nil, errors.New("password length cannot be negative")
	}

	if length == 0 {
		length = defaultPasswordLength
	}

	passwords := make(Passwords, numPasswords)

	// Combined character set for password generation
	charSet := lowercaseLetters + uppercaseLetters + digits + specialChars
	charSetLength := big.NewInt(int64(len(charSet)))

	for i := 0; i < numPasswords; i++ {
		password, err := generatePassword(length, charSet, charSetLength)
		if err != nil {
			return nil, err
		}
		passwords[i] = password
	}

	return passwords, nil
}

// generatePassword creates a single random password
func generatePassword(length int, charSet string, charSetLength *big.Int) (string, error) {
	password := make([]byte, length)

	for i := range password {
		// Get random index from character set
		idx, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			return "", err
		}
		password[i] = charSet[idx.Int64()]
	}

	return string(password), nil
}

// Alternative version with configurable character sets
func GenerateRandomPasswordsAdvanced(length int, useLower, useUpper, useDigits, useSpecial bool) (Passwords, error) {
	if length <= 0 {
		return nil, errors.New("password length must be positive")
	}

	// Build character set based on options
	var charSet string
	if useLower {
		charSet += lowercaseLetters
	}
	if useUpper {
		charSet += uppercaseLetters
	}
	if useDigits {
		charSet += digits
	}
	if useSpecial {
		charSet += specialChars
	}

	if charSet == "" {
		return nil, errors.New("at least one character set must be selected")
	}

	charSetLength := big.NewInt(int64(len(charSet)))
	passwords := make(Passwords, numPasswords)

	for i := 0; i < numPasswords; i++ {
		password, err := generatePassword(length, charSet, charSetLength)
		if err != nil {
			return nil, err
		}
		passwords[i] = password
	}

	return passwords, nil
}

func Hash(password Password, ssid string) []byte {
	hash := pbkdf2.Key([]byte(password), []byte(ssid), 4096, 32, sha1.New)
	return hash
}
