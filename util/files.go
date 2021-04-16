package util

import (
	"math/rand"
	"path/filepath"
	"strings"
)

// SplitName splits the filename info extension and name
func SplitName(filename string) (string, string) {
	ext := filepath.Ext(filename) // Get the extension of the file
	name := strings.TrimSuffix(filename, ext) // Remoe the extension from the original name
	return name, ext // Return both name and extension
}

// GenerateName generates a random name for a file
func GenerateName() string {
	r := make([]byte, 30) // Make a 30 lette array
	// The list of the letters to pick from
	c := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// For each letter pick a random one from the list
	for i := range r {
		r[i] = c[rand.Intn(len(c))]
	}
	// Return as a string
	return string(r)
}

// ToGlobalPath converts file name to filepath in file server
func ToGlobalPath(p string) string {
	return "files/" + p
}

// ToLocalPath converts file name to filepath in file server
func ToLocalPath(p string) string {
	return p[6:]
}
