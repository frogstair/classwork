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
	util.PassClean(&r.Password)
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
		FirstName string
		LastName  string
		ID        string
	}{user.FirstName, user.LastName, user.ID}

	resp.Data = userResponse
	resp.Error = ""

	return 200, resp
}

// LoginUser is the model to create a token for the user
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginUser) clean() {
	util.Clean(&l.Email)
	util.PassClean(&l.Password)
}

// Login will generate a token for the user
func (l *LoginUser) Login(db *gorm.DB) (int, *Response) {
	resp := new(Response)

	user := new(User)
	err := db.Where("email = ?", l.Email).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid email or password"
			return 401, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	if !util.Compare(user.Password, l.Password) {
		resp.Data = nil
		resp.Error = "Invalid email or password"
		return 500, resp
	}

	user.Token = CreateToken(user.ID)
	err = db.Save(user).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	loginResponse := struct {
		Token     string `json:"token"`
		ExpiresIn int64  `json:"expires_in"`
	}{user.Token, 2592000}

	resp.Data = loginResponse
	resp.Error = ""

	return 200, resp
}
