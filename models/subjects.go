package models

import (
	"classwork/util"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// Subject is the internal model of a subject
type Subject struct {
	ID          string        `gorm:"primaryKey" json:"id"`
	TeacherID   string        `gorm:"not null" json:"teacher_id"`
	Teacher     *User         `gorm:"not null" json:"teacher"`
	SchoolID    string        `gorm:"not null" json:"school_id"`
	Name        string        `gorm:"not null" json:"name"`
	Students    []*User       `gorm:"many2many:subject_students" json:"students,omitempty"`
	NumStudents int           `json:"num_students"`
	Assignments []*Assignment `gorm:"many2many:subject_assignments" json:"assignments,omitempty"`
}

// Delete deletes a subject
func (s *Subject) Delete(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

	err := db.Delete(s).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	resp.Data = true
	resp.Error = ""

	return 202, resp
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
func (n *NewSubject) Add(db *gorm.DB, u *User) (int, *util.Response) {
	resp := new(util.Response)

	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	school := new(School)
	err := db.Where("id = ?", n.SchoolID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school id"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	subj := new(Subject)
	subj.ID = ksuid.New().String()
	subj.TeacherID = u.ID
	subj.Teacher = u
	subj.Name = n.Name
	subj.NumStudents = 0
	subj.SchoolID = n.SchoolID

	err = db.Create(subj).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	school.Subjects = append(school.Subjects, subj)
	err = db.Save(school).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	subjResponse := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{subj.ID, subj.Name}

	resp.Data = subjResponse
	resp.Error = ""
	return 201, resp
}

// DeleteSubject deletes a subject
type DeleteSubject struct {
	ID string `json:"id"`
}

func (d *DeleteSubject) clean() {
	util.RemoveSpaces(&d.ID)
}

// Delete deletes a subject
func (d *DeleteSubject) Delete(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response)

	d.clean()

	subject := new(Subject)
	err := db.Where("id = ?", d.ID).First(subject).Error
	if err != nil {
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Invalid subject id"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	if subject.Teacher.ID != user.ID {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	school := new(School)
	err = db.Where("id = ?", school.ID).First(school).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	if school.UserID != user.ID && user.Has(Headmaster) {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	db.Delete(subject)

	resp.Data = true
	resp.Error = ""
	return 202, resp
}

// NewSubjectStudent is the model to add a new student to the subject
type NewSubjectStudent struct {
	ID      string `json:"user_id"`
	Subject string `json:"subject_id"`
}

func (n *NewSubjectStudent) clean() {
	util.RemoveSpaces(&n.ID)
	util.RemoveSpaces(&n.Subject)
}

// Add adds a new student to a subject
func (n *NewSubjectStudent) Add(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response)
	n.clean()

	usr := new(User)
	err := db.Where("id = ?", n.ID).First(usr).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid user id"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	subject := new(Subject)
	err = db.Where("id = ?", n.Subject).First(subject).Error
	if err != nil {
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Invalid subject id"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	school := new(School)
	err = db.Where("id = ?", subject.SchoolID).First(school).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	err = db.Model(school).Association("Students").Find(&school.Students).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	err = db.Model(subject).Association("Students").Find(&subject.Students).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	err = db.Model(user).Association("Subjects").Find(&user.Subjects).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	if subject.TeacherID != user.ID {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	if !usr.Has(Student) {
		resp.Data = nil
		resp.Error = "user not a student"
		return 400, resp
	}

	found := false
	for _, student := range school.Students {
		if student.ID == usr.ID {
			found = true
			break
		}
	}
	if !found {
		resp.Data = nil
		resp.Error = "user not in school"
		return 400, resp
	}

	found = false
	for _, student := range subject.Students {
		if student.ID == usr.ID {
			found = true
			break
		}
	}
	if found {
		resp.Data = nil
		resp.Error = "Student already added"
		return 409, resp
	}

	subject.Students = append(subject.Students, usr)
	subject.NumStudents++
	err = db.Save(subject).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	user.Subjects = append(user.Subjects, subject)
	err = db.Save(user).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	studentResponse := struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{usr.ID, usr.FirstName, usr.LastName}

	resp.Data = studentResponse
	resp.Error = ""

	return 201, resp
}
