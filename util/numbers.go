package util

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

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

// RandomCode generates a random string
func RandomCode() string {
	r := make([]byte, 15) // Make a 15 letter random code
	c := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789$%{}#!@_"
	// Fill in an array with random letters
	for i := range r {
		r[i] = c[rand.Intn(len(c))]
	}
	// Return the hash of that for more randomness
	return Hash(string(r))
}
