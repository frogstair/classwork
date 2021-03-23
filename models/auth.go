package models

import (
	"classwork/util"

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
	if len(r.Password) < 7 {
		return false, "Password must be at least 7 characters"
	}

	return true, ""
}

// Register registers the user intp the database
func (r *RegisterUser) Register(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response) // Response placeholder

	r.clean()                     // Remove trailing whitespace and invalid characters
	valid, reason := r.validate() // Validate the input, and if invalid get reason
	if !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	hashed := util.Hash(r.Password) // Hash the user's password
	user := new(User)               // Create a user placeholder

	user.ID = ksuid.New().String() // Create a GUID and fill in user information
	user.Email = r.Email
	user.FirstName = r.FirstName
	user.LastName = r.LastName
	user.Password = hashed
	user.Perms = Headmaster
	user.PassSet = true

	err := db.Create(user).Error // Create the user in the database
	if err != nil {              // If an error occured
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Email is taken"
			return 409, resp
		}
		return util.DatabaseError(err, resp)
	}

	userResponse := struct { // Create a response for the user
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ID        string `json:"id"`
	}{user.FirstName, user.LastName, user.ID}

	resp.Data = userResponse
	resp.Error = ""

	return 201, resp // Respond
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
func (l *LoginUser) Login(db *gorm.DB) (int, *util.Response, string) {
	resp := new(util.Response) // Create response placeholder
	l.clean()                  // Clean the user input (remove trailing spaces and invalid characters)

	user := new(User)
	err := db.Where("email = ?", l.Email).First(user).Error // Get a user from the database with the entered email

	if err != nil { // If an error occured
		if util.IsNotFoundErr(err) { // If no user was found
			resp.Data = nil
			resp.Error = "Invalid email or password"
			return 401, resp, ""
		}
		_, resp = util.DatabaseError(err, resp) // Any other errors
		return 500, resp, ""
	}

	if user.PassSet { // If user already has a password
		if !util.Compare(user.Password, l.Password) { // Check if password is correct
			resp.Data = nil
			resp.Error = "Invalid email or password"
			return 403, resp, ""
		}
	} else { // Check the code
		if l.Code == user.OneTimeCode { // Create the password if the code is correct
			user.PassSet = true
			user.Password = util.Hash(l.Password)
			user.OneTimeCode = ""
		} else { // Otherwise return nothing
			resp.Data = nil
			resp.Error = "Invatid OTC"
			return 401, resp, ""
		}
	}

	user.Token = CreateToken(user.ID) // Create a token for the user
	err = db.Save(user).Error         // Save the user with the token and password
	if err != nil {
		_, resp = util.DatabaseError(err, resp)
		return 500, resp, ""
	}

	resp.Data = true // Place success flag
	resp.Error = ""

	return 200, resp, user.Token // Respond with the token
}

// OTCCreate is the struct to check is a user has a password
type OTCCreate struct {
	Email string
}

// Create creates an OTC for the user
func (o *OTCCreate) Create(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response) // Response placeholder
	user := new(User)          // User placeholder

	err := db.Where("email = ?", o.Email).First(user).Error // Get the user

	if err != nil { // Check for errors
		if util.IsNotFoundErr(err) { // If user wasn't found
			resp.Data = nil
			resp.Error = "Invalid user"
			return 404, resp
		}
		return util.DatabaseError(err, resp) // If any other error
	}

	if user.PassSet { // If the user has already tried to generate the code
		resp.Data = nil
		resp.Error = "resource gone"
		return 410, resp
	}

	onetimecode := util.RandomCode() // Generate a random code

	user.OneTimeCode = onetimecode // Save the code until the user uses it
	err = db.Save(user).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	resp.Data = onetimecode // Return the code to the user
	resp.Error = ""

	return 201, resp
}
