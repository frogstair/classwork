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

// TokenValidity specifies how long a token is valid
var TokenValidity = int64(2592000)

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
	token.ExpiresAt = time.Now().Unix() + TokenValidity

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token)
	tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}

// ParseToken parses a token and returns the associated user
func ParseToken(tokstr string, db *gorm.DB) (int, *util.Response) {

	resp := new(util.Response)

	tok := &Token{}
	token, err := jwt.ParseWithClaims(tokstr, tok,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

	if !token.Valid {
		if err != nil {
			resp.Data = nil
			resp.Error = "token is invalid"

			return 401, resp
		}
	}

	user := new(User)
	err = db.Where("id = ?", tok.ID).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "token does not correspond to user"

			return 401, resp
		}
		log.Printf("Database error: %s\n", err.Error())
		return util.DatabaseError(err, resp)
	}

	if user.Token != tokstr {
		resp.Data = nil
		resp.Error = "token does not correspond to user"

		return 401, resp
	}

	resp.Data = user
	resp.Error = ""
	return 200, resp
}
