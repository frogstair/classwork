package models

import (
	"classwork/util"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// RegisterUser is the model to register a new user
type RegisterUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *RegisterUser) clean() {
	util.Clean(&r.FirstName)
	util.Clean(&r.LastName)
	util.RemoveSpaces(&r.Password)
}

func (r *RegisterUser) validate() (bool, string) {

	if !util.ValidateName(r.FirstName) {
		return false, "First name should be at least 4 characters"
	}
	if !util.ValidateName(r.LastName) {
		return false, "Last name should be at least 4 characters"
	}
	if !util.ValidateEmail(r.Email) {
		return false, "Email is invalid"
	}

	return true, ""
}

// Register registers the user intp the database
func (r *RegisterUser) Register(db *gorm.DB) (int, *Response) {
	resp := new(Response)

	r.clean()
	valid, reason := r.validate()
	if !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	hashed := util.Hash(r.Password)
	user := new(User)

	user.ID = ksuid.New().String()
	user.Email = r.Email
	user.FirstName = r.FirstName
	user.LastName = r.LastName
	user.Password = hashed
	user.Perms = Headmaster
	user.PassSet = true

	err := db.Create(user).Error
	if err != nil {
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Email is taken"
			return 409, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	userResponse := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ID        string `json:"id"`
	}{user.FirstName, user.LastName, user.ID}

	resp.Data = userResponse
	resp.Error = ""

	return 200, resp
}

// LoginUser is the model to create a token for the user
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func (l *LoginUser) clean() {
	util.RemoveSpaces(&l.Email)
	util.RemoveSpaces(&l.Password)
	util.RemoveSpaces(&l.Code)
}

// Login will generate a token for the user
func (l *LoginUser) Login(db *gorm.DB) (int, *Response, string) {
	resp := new(Response)
	l.clean()

	user := new(User)
	err := db.Where("email = ?", l.Email).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid email or password"
			return 401, resp, ""
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp, ""
	}

	if user.PassSet {
		if !util.Compare(user.Password, l.Password) {
			resp.Data = nil
			resp.Error = "Invalid email or password"
			return 500, resp, ""
		}
	} else {
		if l.Code == user.OneTimeCode {
			user.PassSet = true
			user.Password = util.Hash(l.Password)
			user.OneTimeCode = ""
		} else {
			resp.Data = nil
			resp.Error = "Invatid OTC"
			return 401, resp, ""
		}
	}

	user.Token = CreateToken(user.ID)
	err = db.Save(user).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp, ""
	}

	loginResponse := struct {
		Token     string `json:"token"`
		ExpiresIn int64  `json:"expires_in"`
	}{user.Token, TokenValidity}

	resp.Data = loginResponse
	resp.Error = ""

	return 200, resp, user.Token
}

// OTCCreate is the struct to check is a user has a password
type OTCCreate struct {
	Email string
}

func (o *OTCCreate) clean() {
	util.RemoveSpaces(&o.Email)
}

// Create creates an OTC for the user
func (o *OTCCreate) Create(db *gorm.DB) (int, *Response) {
	resp := new(Response)
	user := new(User)

	err := db.Where("email = ?", o.Email).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid user"
			return 403, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	if user.PassSet {
		resp.Data = nil
		resp.Error = "not found"
		return 404, resp
	}

	onetimecode := util.RandomCode()

	user.OneTimeCode = onetimecode
	err = db.Save(user).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = onetimecode
	resp.Error = ""

	return 200, resp
}
