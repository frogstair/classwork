package models

import (
	"classwork/util"
	"log"

	"github.com/jinzhu/gorm"
)

// User is the internal representation of a user
type User struct {
	ID          string     `gorm:"primaryKey" json:"id"`
	FirstName   string     `gorm:"not null" json:"first_name"`
	LastName    string     `gorm:"not null" json:"last_name"`
	Email       string     `gorm:"not null;unique" json:"email"`
	Password    string     `json:"-"`
	Token       string     `gorm:"unique" json:"-"`
	Perms       Role       `gorm:"not null" json:"-"`
	PassSet     bool       `json:"-"`
	OneTimeCode string     `json:"-"`
	Subjects    []*Subject `gorm:"many2many:subject_students" json:"students,omitempty"`
}

// Has returns if a user has a role
func (u *User) Has(r Role) bool {
	role := u.Perms & r
	return role == r
}

// GetDashboard gets the users dashboard
func (u *User) GetDashboard(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

	dashboard := new(Dashboard)

	if u.Has(Headmaster) {
		hmDashboard := new(HeadmasterDashboard)

		err := db.Where("user_id = ?", u.ID).Find(&hmDashboard.Schools).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}

		dashboard.Headmaster = hmDashboard
	}
	if u.Has(Teacher) {
		tchDashboard := new(TeacherDashboard)

		err := db.Where("teacher_id = ?", u.ID).Find(&tchDashboard.Subjects).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}

		dashboard.Teacher = tchDashboard
	}
	if u.Has(Student) {
		stuDashboard := new(StudentDashboard)

		db.Model(u).Association("Subjects").Find(&u.Subjects)
		for s, subject := range u.Subjects {
			db.Where("subject_id = ?", subject.ID).Find(&subject.Assignments)
			u.Subjects[s] = subject
			for a, assignment := range subject.Assignments {
				db.Model(assignment).Association("Requests").Find(&assignment.Requests)
				u.Subjects[s].Assignments[a] = assignment
				for r, req := range assignment.Requests {
					upl := make([]*RequestUpload, 0)
					req.Uploads = &upl
					db.Model(req).Association("Uploads").Find(&req.Uploads)
					found := false
					for _, upload := range *req.Uploads {
						if upload.UserID == u.ID {
							found = true
							break
						}
					}

					req.Complete = &found
					req.Uploads = nil
					u.Subjects[s].Assignments[a].Requests[r] = req
				}
			}
		}

		stuDashboard.Subject = u.Subjects

		dashboard.Student = stuDashboard
	}

	dashboard.User = u

	resp.Data = dashboard
	resp.Error = ""

	return 200, resp
}

// Delete deletes a user
func (u *User) Delete(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

	err := db.Delete(u).Error
	if err != nil {
		resp.Error = "internal error"
		resp.Data = nil
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = true
	resp.Error = ""

	return 202, resp
}

// Logout removes a token from a user
func (u *User) Logout(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response)

	u.Token = ""
	err := db.Save(u).Error
	if err != nil {
		resp.Error = "internal error"
		resp.Data = nil
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = true
	resp.Error = ""

	return 202, resp
}
