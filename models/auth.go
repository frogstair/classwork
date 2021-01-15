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
	util.RemoveSpaces(&l.Email)
	util.RemoveSpaces(&l.Password)
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

	type loginResponse struct {
		LoggedIn  *bool  `json:"logged_in"`
		OTC       string `json:"otc,omitempty"`
		Token     string `json:"token,omitempty"`
		ExpiresIn *int64 `json:"expires_in,omitempty"`
	}

	loginResp := loginResponse{}

	if user.PassSet {
		if !util.Compare(user.Password, l.Password) {
			resp.Data = nil
			resp.Error = "Invalid email or password"
			return 500, resp, ""
		}

		user.Token = CreateToken(user.ID)
		err = db.Save(user).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			log.Printf("Database error: %s\n", err.Error())
			return 500, resp, ""
		}

		loginResp = loginResponse{&user.PassSet, "", user.Token, &TokenValidity}
	} else {
		onetimecode := util.RandomCode()

		user.OneTimeCode = onetimecode
		err := db.Save(user).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			log.Printf("Database error: %s\n", err.Error())
			return 500, resp, ""
		}

		loginResp = loginResponse{&user.PassSet, onetimecode, "", nil}
	}

	resp.Data = loginResp
	resp.Error = ""

	return 200, resp, user.Token
}

// PasswordCreate is the struct to create a new password for a user
type PasswordCreate struct {
	Password    string
	OneTimeCode string
}

func (p *PasswordCreate) clean() {
	util.RemoveSpaces(&p.Password)
	util.RemoveSpaces(&p.OneTimeCode)
}

// Create creates a password for a given email
func (p *PasswordCreate) Create(db *gorm.DB) (int, *Response, string) {
	resp := new(Response)

	p.clean()

	user := new(User)
	err := db.Where("onetimecode = ?", p.OneTimeCode).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid code"
			return 401, resp, ""
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp, ""
	}

	hashed := util.Hash(p.Password)
	user.Password = hashed
	user.OneTimeCode = ""
	user.PassSet = true
	err = db.Save(user).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp, ""
	}

	l := new(LoginUser)
	l.Email = user.Email
	l.Password = p.Password

	return l.Login(db)
}
