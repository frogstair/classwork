package util

import (
	"math/rand"
	"path/filepath"
	"strings"
)

// SplitName splits the filename info extension and name
func SplitName(filename string) (string, string) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return name, ext
}

// GenerateName generates a random name for a file
func GenerateName() string {
	r := make([]byte, 30)
	c := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	for i := range r {
		r[i] = c[rand.Intn(len(c))]
	}
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
