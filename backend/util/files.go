package util

import (
	"path/filepath"
	"strings"
)

// ToRelativeFPath converts file name to filepath in file server
func ToRelativeFPath(p string) string {
	return "./files/" + p
}

// SplitName splits the filename info extension and name
func SplitName(filename string) (string, string) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return name, ext
}
