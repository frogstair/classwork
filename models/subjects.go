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
	Teacher     *User         `gorm:"not null" json:"teacher,omitempty"`
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
	resp := new(util.Response) // Placeholder response

	n.clean() // Remove trailing whitespace

	if valid, reason := n.validate(); !valid { // Valiate the input
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	school := new(School) // Get the school from database by ID provided
	err := db.Where("id = ?", n.SchoolID).First(school).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid school id"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	db.Model(school).Association("Teachers").Find(&school.Teachers) // Get all teachers from the school

	found := false // Check if the teacher is in the school
	for _, t := range school.Teachers {
		if u.ID == t.ID {
			found = true
			break
		}
	}

	if !found { // If not in school
		resp.Data = nil
		resp.Error = "Teacher not in school"
		return 403, resp
	}

	subj := new(Subject) // Create new subject
	subj.ID = ksuid.New().String()
	subj.TeacherID = u.ID
	subj.Teacher = u
	subj.Name = n.Name
	subj.NumStudents = 0
	subj.SchoolID = n.SchoolID

	err = db.Create(subj).Error // Save it to the database
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	school.Subjects = append(school.Subjects, subj) // Add the subject to the school
	err = db.Save(school).Error // Save the school
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	subjResponse := struct { // Create a response
		ID   string `json:"id"`
		Name string `json:"name"`
	}{subj.ID, subj.Name}

	resp.Data = subjResponse // Respond to the user
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
	resp := new(util.Response) // Placeholder response

	d.clean()

	subject := new(Subject) // Get subject from ID
	err := db.Where("id = ?", d.ID).First(subject).Error
	if err != nil {
		if util.IsDuplicateErr(err) {
			resp.Data = nil
			resp.Error = "Invalid subject id"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	subject.Teacher = new(User) // Get the subject teacher
	db.Model(subject).Association("Teacher").Find(&subject.Teacher)

	if subject.Teacher.ID != user.ID { // If the teacher doesnt own the subject
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	// Headmaster can also delete the subject, so get the headmaster from the school
	school := new(School)
	err = db.Where("id = ?", subject.SchoolID).First(school).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	// Check if the user is a headmaster but doesnt own the school
	if school.UserID != user.ID && user.Has(Headmaster) {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	db.Delete(subject) // Delete the subject from the database

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
		if util.IsNotFoundErr(err) {
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

	usr.Subjects = append(usr.Subjects, subject)
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

// GetSubjectInfo is the internal model to retreive a subject's info
type GetSubjectInfo struct {
	ID  string
}

// Get gets info about the subject
func (g *GetSubjectInfo) Get(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response

	subject := new(Subject) // Get the subject from the database
	err := db.Where("id = ?", g.ID).First(subject).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Unknown subject ID"
			return 404, resp
		}
		return util.DatabaseError(err, resp)
	}

	if user.Has(Teacher) { // If the user is a teacher
		if subject.TeacherID != user.ID { // If the teacher doesnt own the subjet
			resp.Data = nil
			resp.Error = "forbidden"
			return 403, resp
		}

		db.Model(subject).Association("Students").Find(&subject.Students) // Get all the students for the subject

		// Get ten latest assignments for the subject
		db.Where("subject_id = ?", g.ID).Order("time_assigned desc").Limit(10).Find(&subject.Assignments)
		for a, assignment := range subject.Assignments {
			// For each assignment get the requests
			db.Model(assignment).Association("Requests").Find(&assignment.Requests)
			// For each request
			for _, req := range assignment.Requests {
				// Get all student uploads
				upl := make([]*RequestUpload, 0)
				req.Uploads = upl
				db.Model(req).Association("Uploads").Find(&req.Uploads)
				req.Complete = nil
			}
			
			// Get the list of students that completed it
			db.Model(assignment).Association("CompletedBy").Find(&assignment.CompletedBy)

			// If the assignment has a due time
			if assignment.TimeDue != nil {
				// Get list of completed users
				completed := make(map[string]bool)
				// For each user that completed the assignment make a record of that
				for _, compl := range assignment.CompletedBy {
					completed[compl.ID] = true
				}

				// make an array of users that didnt complete it
				assignment.NotCompletedBy = make([]*User, 0)
				// Fill out the list of students that didnt complete it
				for _, student := range subject.Students {
					if _, ok := completed[student.ID]; !ok {
						assignment.NotCompletedBy = append(assignment.NotCompletedBy, student)
					}
				}
			}

			subject.Assignments[a] = assignment
		}

		// Get all students from the subject
		getStudents := new(GetStudents)
		getStudents.ID = g.SID
		code, students := getStudents.Get(db)

		// If an error occured while getting the user,
		// respond with its error code instead
		if code != 200 {
			return code, students
		}

		// Respond
		response := struct {
			S  *Subject `json:"subject"`
			St []*User  `json:"students"`
		}{
			subject,
			students.Data.([]*User),
		}

		resp.Data = response
		resp.Error = ""

	} else { // If the user isnt a teacher then forbid access
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	return 200, resp
}
