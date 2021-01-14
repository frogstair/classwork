package util

import "golang.org/x/crypto/bcrypt"

// Hash will hash a string and return it
func Hash(hash string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(hash), 10)
	return string(hashed)
}

// Compare compares a hashed password to a plaintext password
func Compare(hash string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	return err == nil
}
