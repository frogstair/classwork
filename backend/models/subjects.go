package models

import (
	"classwork/backend/util"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// Subject is the internal model of a subject
type Subject struct {
	ID          string  `gorm:"primaryKey" json:"id"`
	UserID      string  `gorm:"not null" json:"teacher_id"`
	SchoolID    string  `gorm:"not null" json:"school_id"`
	Name        string  `gorm:"not null" json:"name"`
	Students    []*User `gorm:"many2many:students" json:"students,omitempty"`
	NumStudents int     `json:"num_students"`
}

// NewSubject is a new subject
type NewSubject struct {
	Name     string `json:"name"`
	SchoolID string `json:"school_id"`
}

func (n *NewSubject) clean() {
	util.Clean(&n.Name)
	util.RemoveSpaces(&n.SchoolID)
}

func (n *NewSubject) validate() (bool, string) {
	if !util.ValidateName(n.Name) {
		return false, "Name should be at least four characters"
	}
	return true, ""
}

// Add adds a subject to the database
func (n *NewSubject) Add(db *gorm.DB, u *User) (int, *Response) {
	resp := new(Response)

	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	school := new(School)
	err := db.Where("id = ?", n.SchoolID).First(school).Error
	if err != nil {
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school id"
			return 400, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	subj := new(Subject)
	subj.ID = ksuid.New().String()
	subj.UserID = u.ID
	subj.Name = n.Name
	subj.NumStudents = 0
	subj.SchoolID = n.SchoolID

	err = db.Create(subj).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	subjResponse := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{subj.ID, subj.Name}

	resp.Data = subjResponse
	resp.Error = ""
	return 200, resp
}

// DeleteSubject deletes a subject
type DeleteSubject struct {
	ID string `json:"id"`
}

func (d *DeleteSubject) clean() {
	util.RemoveSpaces(&d.ID)
}

// Delete deletes a subject
func (d *DeleteSubject) Delete(db *gorm.DB, user *User) (int, *Response) {
	resp := new(Response)

	d.clean()

	subject := new(Subject)
	err := db.Where("id = ?", d.ID).First(subject).Error
	if err != nil {
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Invalid subject id"
			return 400, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	if subject.UserID != user.ID {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	school := new(School)
	err = db.Where("id = ?", school.ID).First(school).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	if school.UserID != user.ID && user.Has(Headmaster) {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	db.Delete(subject)

	resp.Data = true
	resp.Error = ""
	return 200, resp
}
