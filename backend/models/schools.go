package models

import (
	"classwork/backend/util"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// School is the internal representation of the schools
type School struct {
	ID       string  `gorm:"primaryKey" json:"id"`
	UserID   string  `gorm:"not null" json:"-"`
	Name     string  `gorm:"not null" json:"name"`
	Students []*User `gorm:"many2many:students" json:"students,omitempty"`
	Teachers []*User `gorm:"many2many:teachers" json:"teachers,omitempty"`
}

// NewSchool is the model to add a new school
type NewSchool struct {
	Name string `json:"name"`
}

func (n *NewSchool) clean() {
	util.Clean(&n.Name)
}

func (n *NewSchool) validate() (bool, string) {
	if !util.ValidateName(n.Name) {
		return false, "Name should be at least 4 characters"
	}
	return true, ""
}

// Add adds a new school to the database
func (n *NewSchool) Add(db *gorm.DB, user *User) (int, *Response) {
	resp := new(Response)

	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Error = reason
		resp.Data = nil
		return 400, resp
	}

	school := new(School)

	school.Name = n.Name
	school.UserID = user.ID
	school.ID = ksuid.New().String()

	err := db.Save(school).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	schoolResp := struct {
		Name string
		ID   string
	}{n.Name, school.ID}

	resp.Data = schoolResp
	resp.Error = ""
	return 200, resp
}

// DeleteSchool deletes a school
type DeleteSchool struct {
	ID string `json:"id"`
}

func (d *DeleteSchool) clean() {
	util.RemoveSpaces(&d.ID)
}

// Delete will delete a school
func (d *DeleteSchool) Delete(db *gorm.DB, user *User) (int, *Response) {
	resp := new(Response)
	d.clean()
	school := new(School)

	err := db.Where("id = ?", d.ID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school ID"
			return 400, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	if school.UserID != user.ID {
		resp.Data = nil
		resp.Error = "user does not own school"
		return 403, resp
	}

	err = db.Delete(school).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = true
	resp.Error = ""

	return 200, resp
}

// GetSchoolInfo is the model to get school info
type GetSchoolInfo struct {
	ID string `json:"id"`
}

func (g *GetSchoolInfo) clean() {
	util.RemoveSpaces(&g.ID)
}

// GetInfo gets the info for a school
func (g *GetSchoolInfo) GetInfo(db *gorm.DB, user *User) (int, *Response) {

	resp := new(Response)

	school := new(School)
	db.Where("id = ?", g.ID).First(school)

	if school.UserID != user.ID {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	db.Model(school).Association("Teachers").Find(&school.Teachers)
	db.Model(school).Association("Students").Find(&school.Students)

	found := false
	for _, teacher := range school.Teachers {
		if teacher.ID == user.ID {
			found = true
			break
		}
	}
	for _, student := range school.Students {
		if found {
			break
		}
		if student.ID == user.ID {
			found = true
			break
		}
	}

	if !found {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	if user.Has(Headmaster) {

	} else if user.Has(Teacher) {
		school.Teachers = make([]*User, 0)
		school.Students = make([]*User, 0)
	} else if user.Has(Student) {
		school.Teachers = make([]*User, 0)
		school.Students = make([]*User, 0)
	}

	resp.Data = school

	return 200, resp
}
