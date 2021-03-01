package models

import (
	"classwork/util"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Role is a type alias for a role a user can have
type Role int

// Role list
const (
	Headmaster Role = 1 << iota // 001
	Teacher                     // 010
	Student                     // 100
)

// TokenValidity specifies how long a token is valid
var TokenValidity = int64(2592000) // 30 days

// Token is the model for a JWT token with custom claims
type Token struct {
	jwt.StandardClaims // Standard JWT claims (in my case only the `expires` field)

	ID string `json:"id"` // The user's ID
}

// CreateToken will create a token string for a given ID
func CreateToken(id string) string {
	token := new(Token) // Create a new token

	token.ID = id // Set the ID in the token to the passed ID

	token.ExpiresAt = time.Now().Unix() + TokenValidity // Set expiration time (current time + how long a token is valid)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token)             // Create a new JWT token instance signed with HS512
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET"))) // Encrypt the token with the secret from an environment variable

	return tokenString // Return the string to be sent to the user
}

// ParseToken parses a token and returns the associated user
func ParseToken(tokstr string, db *gorm.DB) (int, *util.Response) {

	resp := new(util.Response) // Create a response to the user

	tok := &Token{}                                // Allocate memory for a new token
	token, err := jwt.ParseWithClaims(tokstr, tok, // ParseWithClaims function parses the token with a function that returns the key
		func(t *jwt.Token) (interface{}, error) { // Since I only store the key in an environment variable, get the key from there
			return []byte(os.Getenv("JWT_SECRET")), nil // Just return it
		})

	if !token.Valid { // If the token is invalid (i.e invalid syntax)
		if err != nil {
			resp.Data = nil
			resp.Error = "token is invalid"

			return 401, resp // Respond with corresponding message and error code
		}
	}

	user := new(User) // Retreive a user from the database
	err = db.Where("id = ?", tok.ID).First(user).Error
	if err != nil { // Error handling
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "token does not correspond to user"

			return 401, resp
		}
		log.Printf("Database error: %s\n", err.Error())
		return util.DatabaseError(err, resp)
	}

	if user.Token != tokstr { // In the end, if the tokens dont match
		resp.Data = nil
		resp.Error = "token does not correspond to user"

		return 401, resp
	}

	resp.Data = user // Return the user variable
	resp.Error = ""
	return 200, resp // Success
}
