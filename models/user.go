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
	Token       string     `json:"-"`
	Perms       Role       `gorm:"not null" json:"-"`
	PassSet     bool       `json:"-"`
	OneTimeCode string     `json:"-"`
	Subjects    []*Subject `gorm:"many2many:subject_students" json:"students,omitempty"`
}

// Has returns if a user has a role
func (u *User) Has(r Role) bool {
	role := u.Perms & r // Perform AND logical operation
	return role == r    // If the AND operation did not return 0, return true
}

// GetDashboard gets the users dashboard
func (u *User) GetDashboard(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response

	dashboard := new(Dashboard) // Dashboard placeholder

	if u.Has(Headmaster) { // If the user is a headmaster
		hmDashboard := new(HeadmasterDashboard) // Placeholder

		err := db.Where("user_id = ?", u.ID).Find(&hmDashboard.Schools).Error // Get all the schools the headmaster owns
		if err != nil {
			return util.DatabaseError(err, resp)
		}

		dashboard.Headmaster = hmDashboard // Set the dashboard
	}
	if u.Has(Teacher) { // If the user is a teacher
		tchDashboard := new(TeacherDashboard)

		result := struct { // Result placeholder
			SchoolID  string
			TeacherID string
		}{}
		db.Raw("select * from school_teachers where user_id = ?", u.ID).Scan(&result) // Get all the schools the teacher is in

		tchDashboard.SchoolID = result.SchoolID

		err := db.Where("teacher_id = ?", u.ID).Find(&tchDashboard.Subjects).Error // Get all the teacher subjects
		if err != nil {
			return util.DatabaseError(err, resp)
		}

		dashboard.Teacher = tchDashboard // Set the dashboard
	}
	if u.Has(Student) { // Is user is a student
		stuDashboard := new(StudentDashboard)

		db.Model(u).Association("Subjects").Find(&u.Subjects) // Get for each subject
		for s, subject := range u.Subjects {
			db.Where("subject_id = ?", subject.ID).Order("time_assigned desc").Limit(10).Find(&subject.Assignments)
			u.Subjects[s] = subject
			for a, assignment := range subject.Assignments { // Each assignment
				db.Model(assignment).Association("Requests").Find(&assignment.Requests)
				assignment.CompletedBy = []*User{}
				u.Subjects[s].Assignments[a] = assignment
				for r, req := range assignment.Requests { // Each upload request
					upl := make([]*RequestUpload, 0)
					req.Uploads = upl
					db.Model(req).Association("Uploads").Find(&req.Uploads)
					found := false
					for _, upload := range req.Uploads {
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

	dashboard.User = u // Respond with dashboard

	resp.Data = dashboard
	resp.Error = ""

	return 200, resp
}

// Delete deletes a user
func (u *User) Delete(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response

	err := db.Delete(u).Error // Delete user from database
	if err != nil {
		resp.Error = "internal error"
		resp.Data = nil
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = true // success
	resp.Error = ""

	return 202, resp // Deleted
}

// Logout removes a token from a user
func (u *User) Logout(db *gorm.DB) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response

	u.Token = ""            // Set user token to empty
	err := db.Save(u).Error // Save
	if err != nil {
		resp.Error = "internal error"
		resp.Data = nil
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	resp.Data = true // Success
	resp.Error = ""

	return 202, resp
}
