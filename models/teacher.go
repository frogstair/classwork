package models

import (
	"classwork/util"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// NewTeacher is the model to add a new teacher
type NewTeacher struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	SchoolID  string `json:"school_id"`
}

func (n *NewTeacher) clean() {
	util.RemoveSpaces(&n.Email)
	util.RemoveSpaces(&n.SchoolID)
	util.Clean(&n.FirstName)
	util.Clean(&n.LastName)
}

func (n *NewTeacher) validate() (bool, string) {
	if !util.ValidateEmail(n.Email) {
		return false, "Email is invalid"
	}
	if !util.ValidateName(n.FirstName) {
		return false, "First name should be at least 4 characters"
	}
	if !util.ValidateName(n.LastName) {
		return false, "Last name should be at least 4 characters"
	}
	return true, ""
}

// Add adds a new teacher to the database
func (n *NewTeacher) Add(db *gorm.DB) (int, *Response) {
	resp := new(Response)

	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	user := new(User)
	err := db.Where("email = ?", n.Email).First(user).Error
	found := !util.IsNotFoundErr(err)
	if err != nil && found {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	school := new(School)
	err = db.Where("id = ?", n.SchoolID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "School not found"
			return 404, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}
	if found {
		user.Perms |= Teacher
		school.Teachers = append(school.Teachers, user)

		err = db.Save(user).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			log.Printf("Database error: %s\n", err.Error())
			return 500, resp
		}

		err = db.Save(school).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			log.Printf("Database error: %s\n", err.Error())
			return 500, resp
		}
	} else {
		user = new(User)

		user.ID = ksuid.New().String()
		user.Email = n.Email
		user.FirstName = n.FirstName
		user.LastName = n.LastName
		user.Perms = Teacher
		user.PassSet = false

		err = db.Create(user).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			log.Printf("Database error: %s\n", err.Error())
			return 500, resp
		}
	}

	newTeacherResponse := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ID        string `json:"id"`
	}{user.FirstName, user.LastName, user.ID}

	resp.Data = newTeacherResponse
	resp.Error = ""
	return 200, resp
}
