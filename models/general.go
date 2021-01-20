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
	resp := new(util.Response)
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
		return util.DatabaseError(err, resp)
	}

	resp.Data = false
	resp.Error = ""

	return 200, resp
}
