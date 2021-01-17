package models

import (
	"classwork/backend/util"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// Assignment is the internal model for assignments
type Assignment struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	TeacherID string    `gorm:"not null" json:"teacher_id"`
	Teacher   *User     `gorm:"not null" json:"teacher"`
	SubjectID string    `gorm:"not null" json:"subject_id"`
	Name      string    `gorm:"not null" json:"name"`
	Text      string    `json:"text"`
	Files     string    `json:"files"`
	TimeDue   time.Time `json:"time_due"`
}

// NewAssignment is a model to create a new assignment
type NewAssignment struct {
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	SubjectID string    `json:"subject_id"`
	TimeDue   time.Time `json:"time_due"`
	Files     []string  `json:"files"`
}

func (n *NewAssignment) clean() {
	util.Clean(&n.Name)
	util.Clean(&n.Text)
	util.RemoveSpaces(&n.SubjectID)
	for i := range n.Files {
		util.Clean(&n.Files[i])
	}
}

func (n *NewAssignment) validate() (bool, string) {
	if !util.ValidateName(n.Name) {
		return false, "Name should be at least 4 characters"
	}
	return true, ""
}

// Create creates a new assignment
func (n *NewAssignment) Create(db *gorm.DB, user *User) (int, *Response) {
	resp := new(Response)
	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	if n.TimeDue.Before(time.Now()) {
		resp.Data = nil
		resp.Error = "Cannot set time due in the past"
		return 400, resp
	}

	subject := new(Subject)
	err := db.Where("id = ?", n.SubjectID).Find(subject).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid subject ID"
			return 403, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	if subject.TeacherID != user.ID {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	names := make([]string, len(n.Files))

	for i, file := range n.Files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			resp.Data = nil
			resp.Error = "Internal error"
			return 500, resp
		}

		name, ext := util.SplitName(file)
		name = name[:len(name)-2] + "_1"
		os.Rename(file, name+ext)
		names[i] = name + ext
	}

	assignment := new(Assignment)
	assignment.ID = ksuid.New().String()
	assignment.Name = n.Name
	assignment.Text = n.Text
	assignment.TeacherID = user.ID
	assignment.SubjectID = subject.ID
	assignment.TimeDue = n.TimeDue
	assignment.Files = strings.Join(names, ";")

	err = db.Save(assignment).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	assgn := struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		TeacherID string   `json:"teacher_id"`
		SubjectID string   `json:"subject_id"`
		Files     []string `json:"files"`
	}{assignment.ID, assignment.Name, assignment.TeacherID, assignment.TeacherID, names}

	resp.Data = assgn
	resp.Error = ""

	return 200, resp
}
