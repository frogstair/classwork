package models

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Role is a type alias for a role a user can have
type Role int

var tokenValidity = int64(2592000)

// Role list
const (
	Headmaster Role = 1 << iota
	Teacher
	Student
)

// Token is the model for a JWT token with custom claims
type Token struct {
	jwt.StandardClaims
	ID string `json:"id"`
}

// CreateToken will create a token string for a given ID
func CreateToken(id string) string {
	token := new(Token)

	token.ID = id
	token.ExpiresAt = time.Now().Unix() + tokenValidity

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}
