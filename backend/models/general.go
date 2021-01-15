package models

import (
	"classwork/backend/util"
	"log"

	"github.com/jinzhu/gorm"
)

// Response is the response struct, that will be sent back to the user
type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// Email validates if the email supplied is correct
type Email struct {
	Email string
}

func (e *Email) clean() {
	util.RemoveSpaces(&e.Email)
}

func (e *Email) validate() (bool, string) {
	valid := util.ValidateEmail(e.Email)
	if valid {
		return valid, ""
	}
	return valid, "Invalid email"
}

// Valid returns if the email is valid
func (e *Email) Valid(db *gorm.DB) (int, *Response) {
	resp := new(Response)
	e.clean()
	user := new(User)

	valid, reason := e.validate()

	if !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	err := db.Where("email = ?", e.Email).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = true
			resp.Error = ""
			return 200, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = false
	resp.Error = ""

	return 200, resp
}
