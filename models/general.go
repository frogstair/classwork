package models

import (
	"classwork/util"

	"github.com/jinzhu/gorm"
)

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
func (e *Email) Valid(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response
	e.clean()
	user := new(User)

	valid, reason := e.validate() // Validate the email

	if !valid { // If not valid
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	err := db.Where("email = ?", e.Email).First(user).Error // Get a user by email

	if err != nil { // If error occured
		if util.IsNotFoundErr(err) {
			resp.Data = true // We are actually looking for a not found error
			resp.Error = ""
			return 200, resp
		}
		return util.DatabaseError(err, resp)
	}

	resp.Data = false // If found then return error
	resp.Error = ""

	return 200, resp
}
