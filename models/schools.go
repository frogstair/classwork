package models

import (
	"classwork/util"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// School is the internal representation of the schools
type School struct {
	ID       string     `gorm:"primaryKey" json:"id"`
	UserID   string     `gorm:"not null" json:"-"`
	Name     string     `gorm:"not null" json:"name"`
	Students []*User    `gorm:"many2many:school_students" json:"students,omitempty"`
	Teachers []*User    `gorm:"many2many:school_teachers" json:"teachers,omitempty"`
	Subjects []*Subject `gorm:"many2many:school_subjects" json:"subjects,omitempty"`
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
func (n *NewSchool) Add(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response)

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
		return util.DatabaseError(err, resp)
	}

	schoolResp := struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}{n.Name, school.ID}

	resp.Data = schoolResp
	resp.Error = ""
	return 201, resp
}

// DeleteSchool deletes a school
type DeleteSchool struct {
	ID string `json:"id"`
}

func (d *DeleteSchool) clean() {
	util.RemoveSpaces(&d.ID)
}

// Delete will delete a school
func (d *DeleteSchool) Delete(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response)
	d.clean()
	school := new(School)

	err := db.Where("id = ?", d.ID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school ID"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	if school.UserID != user.ID {
		resp.Data = nil
		resp.Error = "user does not own school"
		return 403, resp
	}

	err = db.Delete(school).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	subjects := make([]*Subject, 0)
	err = db.Where("school_id = ?", school.ID).Find(&subjects).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	for _, subject := range subjects {
		go subject.Delete(db)
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
func (g *GetSchoolInfo) GetInfo(db *gorm.DB, user *User) (int, *util.Response) {

	resp := new(util.Response)

	school := new(School)
	err := db.Where("id = ?", g.ID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school ID"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	if school.UserID != user.ID && user.Has(Headmaster) {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	school.Teachers = make([]*User, 0)
	school.Students = make([]*User, 0)
	school.Subjects = make([]*Subject, 0)

	db.Model(school).Association("Teachers").Find(&school.Teachers)
	db.Model(school).Association("Students").Find(&school.Students)
	db.Model(school).Association("Subjects").Find(&school.Subjects)

	for i, subj := range school.Subjects {
		db.Where("subject_id = ?", subj.ID).Find(&school.Subjects[i].Assignments)
		for j, assignment := range school.Subjects[i].Assignments {
			db.Model(assignment).Association("Requests").Find(&school.Subjects[i].Assignments[j].Requests)
			db.Model(assignment).Association("Files").Find(&school.Subjects[i].Assignments[j].Files)
			db.Model(assignment).Association("CompletedBy").Find(&school.Subjects[i].Assignments[j].CompletedBy)
		}
	}

	for i, subject := range school.Subjects {
		usr := new(User)
		err := db.Where("id = ?", subject.TeacherID).First(usr).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}
		school.Subjects[i].Teacher = usr
	}

	if !user.Has(Headmaster) {
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
	}

	if user.Has(Headmaster) {

	} else if user.Has(Teacher) {
		school.Teachers = make([]*User, 0)
		school.Students = make([]*User, 0)

		subj := make([]*Subject, 0)
		for _, subject := range school.Subjects {
			if subject.TeacherID == user.ID {
				subj = append(subj, subject)
			}
		}

		school.Subjects = subj
	} else if user.Has(Student) {
		school.Teachers = make([]*User, 0)
		school.Students = make([]*User, 0)
		school.Subjects = make([]*Subject, 0)
	}

	resp.Data = school

	return 200, resp
}

// GetStudents is the model to get students from the school
type GetStudents struct {
	ID string
}

func (g *GetStudents) clean() {
	util.RemoveSpaces(&g.ID)
}

// Get gets the students from the school
func (g *GetStudents) Get(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

	g.clean()

	school := new(School)
	err := db.Where("id = ?", g.ID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school ID"
			return 404, resp
		}
		return util.DatabaseError(err, resp)
	}

	db.Model(school).Association("Students").Find(&school.Students)

	if len(school.Students) == 0 {
		resp.Data = []*User{}
	} else {
		resp.Data = school.Students
	}

	resp.Error = ""

	return 200, resp
}
