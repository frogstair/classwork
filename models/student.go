package models

import (
	"classwork/util"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// NewStudent is a model to add a new student
type NewStudent struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	SchoolID  string `json:"school_id"`
}

func (n *NewStudent) clean() {
	util.RemoveSpaces(&n.Email)
	util.RemoveSpaces(&n.SchoolID)
	util.Clean(&n.FirstName)
	util.Clean(&n.LastName)
}

func (n *NewStudent) validate() (bool, string) {
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

// Add adds a new student to the database
func (n *NewStudent) Add(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

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
		return util.DatabaseError(err, resp)
	}

	school := new(School)
	err = db.Where("id = ?", n.SchoolID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "School not found"
			return 404, resp
		}
		return util.DatabaseError(err, resp)
	}

	if found {
		user.Perms |= Student
		school.Students = append(school.Students, user)

		err = db.Save(user).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}

		err = db.Save(school).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}
	} else {
		user = new(User)

		user.ID = ksuid.New().String()
		user.Email = n.Email
		user.FirstName = n.FirstName
		user.LastName = n.LastName
		user.Perms = Student
		user.PassSet = false

		err = db.Create(user).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}

		school.Students = append(school.Students, user)

		err = db.Save(school).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}
	}

	newStudentResponse := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ID        string `json:"id"`
	}{user.FirstName, user.LastName, user.ID}

	resp.Data = newStudentResponse
	resp.Error = ""
	return 201, resp
}

// DeleteStudent is a model to delete a student from a database
type DeleteStudent struct {
	UserID   string `json:"id"`
	SchoolID string `json:"school_id"`
}

func (d *DeleteStudent) clean() {
	util.RemoveSpaces(&d.SchoolID)
	util.RemoveSpaces(&d.UserID)
}

// Delete deletes a teacher
func (d *DeleteStudent) Delete(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

	d.clean()

	user := new(User)
	err := db.Where("id = ?", d.UserID).First(user).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Teacher not found"
			return 404, resp
		}
		return util.DatabaseError(err, resp)
	}

	school := new(School)
	err = db.Where("id = ?", d.SchoolID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Teacher not found"
			return 404, resp
		}
		return util.DatabaseError(err, resp)
	}

	db.Model(school).Association("Students").Delete(user)
	user.Perms &^= Student

	if user.Perms == 0 {
		return user.Delete(db)
	}

	err = db.Save(user).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	resp.Data = true
	resp.Error = ""

	return 202, resp
}
