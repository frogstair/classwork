package models

import (
	"classwork/util"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// Assignment is the internal model for assignments
type Assignment struct {
	ID             string            `gorm:"primaryKey" json:"id"`
	TeacherID      string            `gorm:"not null" json:"teacher_id"`
	Teacher        *User             `gorm:"not null" json:"teacher,omitempty"`
	SubjectID      string            `gorm:"not null" json:"subject_id,omitempty"`
	Name           string            `gorm:"not null" json:"name"`
	Text           string            `json:"text"`
	Files          []*AssignmentFile `json:"files,omitempty"`
	TimeDue        *time.Time        `json:"time_due,omitempty"`
	TimeAssigned   *time.Time        `json:"time_assigned"`
	Requests       []*Request        `json:"requests,omitempty"`
	CompletedBy    []*User           `gorm:"many2many:assignments_completed" json:"comleted_by,omitempty"`
	NotCompletedBy []*User           `gorm:"-" json:"not_completed_by,omitempty"`
}

// AssignmentFile is the internal structure of a file relating to an assignment
type AssignmentFile struct {
	AssignmentID string `json:"-"`
	Path         string `gorm:"primaryKey" json:"path"`
	Name         string `json:"name"`
}

// Request is the model to request an upload for the students
type Request struct {
	ID           string           `gorm:"primaryKey" json:"id"`
	AssignmentID string           `json:"-"`
	Name         string           `gorm:"not null" json:"name"`
	Complete     *bool            `gorm:"-" json:"complete,omitempty"`
	Uploads      []*RequestUpload `json:"uploads,omitempty"`
}

// RequestUpload is the model to tracks who uploaded what
type RequestUpload struct {
	RequestID string `json:"-"`
	UserID    string `json:"-"`
	User      *User  `gorm:"-" json:"user"`
	Filepath  string `json:"path"`
	Filename  string `json:"name"`
}

// NewAssignment is a model to create a new assignment
type NewAssignment struct {
	Name      string `json:"name"`
	Text      string `json:"text"`
	SubjectID string `json:"subject_id"`
	// TimeDue field is a pointer because it is optional
	// if the pointer is nil, we can ignore it and say
	// there is no time due for the request
	TimeDue        *time.Time `json:"time_due"`
	Files          []string   `json:"files"`
	UploadRequest  bool       `json:"-"`
	UploadRequests []string   `json:"uploads"`
}

func (n *NewAssignment) clean() {
	util.Clean(&n.Name)
	util.Clean(&n.Text)
	util.RemoveSpaces(&n.SubjectID)
	for i := range n.Files {
		util.RemoveSpaces(&n.Files[i])
	}
	n.UploadRequest = len(n.UploadRequests) != 0
}

func (n *NewAssignment) validate() (bool, string) {
	if !util.ValidateName(n.Name) {
		return false, "Name should be at least 4 characters"
	}
	return true, ""
}

// Create creates a new assignment
func (n *NewAssignment) Create(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response
	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	// Check if the time due is in the past, and return an error
	if n.UploadRequest && (n.TimeDue == nil || n.TimeDue.Before(time.Now())) {
		resp.Data = nil
		resp.Error = "Cannot set time due in the past"
		return 400, resp
	}

	// Get the subject by ID to which the assignment is added
	subject := new(Subject)
	err := db.Where("id = ?", n.SubjectID).Find(subject).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "Invalid subject ID"
			return 403, resp
		}
		return util.DatabaseError(err, resp)
	}

	// Check if the user owns the subject
	if subject.TeacherID != user.ID {
		resp.Data = nil
		resp.Error = "forbidden"
		return 403, resp
	}

	// Generate filenames for the assignments

	// Create an empty array
	names := make([]string, len(n.Files))
	for i, file := range n.Files {
		// Get the file on disk
		file = util.ToGlobalPath(file)
		// If it doesnt exist then throw an error
		if _, err := os.Stat(file); os.IsNotExist(err) {
			resp.Data = nil
			resp.Error = "Internal error"
			return 500, resp
		}

		// Split the filename and extension
		name, ext := util.SplitName(file)
		// Set file as verified
		name = name[:len(name)-2] + "_1"
		// Rename the file to validate it
		os.Rename(file, name+ext)
		// Place the name in the array
		names[i] = name + ext
	}

	now := time.Now()

	// Set all the fields in the assignment
	assignment := new(Assignment)
	assignment.ID = ksuid.New().String()
	assignment.Name = n.Name
	assignment.Text = n.Text
	assignment.TeacherID = user.ID
	assignment.SubjectID = subject.ID
	assignment.TimeDue = n.TimeDue
	assignment.TimeAssigned = &now

	// Save the assignment in the database
	err = db.Save(assignment).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	// Create all the necessary upload requests
	requests := make([]*Request, len(n.UploadRequests))

	// Generate all the necessary info for all the requests
	for i, uploadReq := range n.UploadRequests {
		request := new(Request)
		request.ID = ksuid.New().String()
		request.AssignmentID = assignment.ID
		request.Name = uploadReq

		requests[i] = request
		db.Save(request)
	}

	// Create all links to the files that are attached to the assignment
	files := make([]*AssignmentFile, len(names))
	for i, name := range names {
		file := new(AssignmentFile)
		file.AssignmentID = assignment.ID
		file.Path = name
		file.Name = util.ToLocalPath(name)

		files[i] = file
		db.Save(file)
	}

	// Make a response
	assgn := struct {
		ID         string            `json:"id"`
		Name       string            `json:"name"`
		TeacherID  string            `json:"teacher_id"`
		SubjectID  string            `json:"subject_id"`
		AssignedAt time.Time         `json:"time_assigned"`
		Files      []*AssignmentFile `json:"files"`
		Uploads    []*Request        `json:"uploads"`
	}{assignment.ID, assignment.Name, assignment.TeacherID, assignment.SubjectID, now, files, requests}

	resp.Data = assgn
	resp.Error = ""

	return 200, resp
}

// NewRequestComplete is a model to complete an upload request
type NewRequestComplete struct {
	RequestID string `json:"request_id"`
	Filename  string `json:"filepath"`
}

func (n *NewRequestComplete) clean() {
	util.RemoveSpaces(&n.Filename)
	util.RemoveSpaces(&n.RequestID)
}

// Complete completes the upload request
func (n *NewRequestComplete) Complete(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response) // Placeholder response

	n.clean() // Remove trailing whitespace

	request := new(Request) // Get the request to complete from the database
	err := db.Where("id = ?", n.RequestID).First(request).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "request not found"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}

	assgn := new(Assignment) // Get the assignment from which the request is
	err = db.Where("id = ?", request.AssignmentID).First(assgn).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	// Fill in all the requests
	db.Model(assgn).Association("Requests").Find(&assgn.Requests)

	subj := new(Subject) // Get the subject from which the assignment is
	err = db.Where("id = ?", assgn.SubjectID).First(subj).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	// Get all the students from the subject
	db.Model(subj).Association("Students").Find(&subj.Students)
	// Check if the student takes the subject
	found := false
	for _, student := range subj.Students {
		if student.ID == user.ID {
			found = true
			break
		}
	}
	if !found {
		resp.Data = nil
		resp.Error = "user not in subject"
		return 403, resp
	}

	// Get all the uplopads from the request to check if a file was already uploaded
	db.Model(request).Association("Uploads").Find(&request.Uploads)
	for _, upl := range request.Uploads {
		if upl.UserID == user.ID {
			resp.Data = nil
			resp.Error = "file already uploaded"
			return 409, resp
		}
	}

	// Check if the file that the user uploaded exists
	if _, err := os.Stat(util.ToGlobalPath(n.Filename)); os.IsNotExist(err) {
		resp.Data = nil
		resp.Error = "file doesnt exist"
		return 400, resp
	}

	// If all checks passed
	// Create a request upload
	reqUpl := new(RequestUpload)
	reqUpl.Filepath = util.ToGlobalPath(n.Filename)
	reqUpl.Filename = n.Filename
	reqUpl.RequestID = request.ID
	reqUpl.UserID = user.ID

	// Save into database
	err = db.Create(reqUpl).Error
	if err != nil {
		return util.DatabaseError(err, resp)
	}

	// Check if the student has any remaining requests unfilled
	completed := true
	db.Model(assgn).Association("Requests").Find(&assgn.Requests)
	db.Model(assgn).Association("CompletedBy").Find(&assgn.CompletedBy)
	for _, req := range assgn.Requests {
		db.Model(req).Association("Uploads").Find(&req.Uploads)
		found := false
		for _, upload := range req.Uploads {
			if upload.UserID == user.ID {
				found = true
				break
			}
		}
		if !found {
			completed = false
			break
		}
	}

	// if all the requests were completed
	if completed {
		// Here I use a raw database query because
		// I need to create an association with existing data
		// which GORM doesnt yet support
		err = db.Exec("insert into assignments_completed(assignment_id, user_id) values (?, ?);", assgn.ID, user.ID).Error
		if err != nil {
			return util.DatabaseError(err, resp)
		}
	}

	resp.Data = true
	resp.Error = ""

	return 202, resp
}

// GetAssignment is the model to get the assignment
type GetAssignment struct {
	ID string
}

// Get gets the assignment information
func (g *GetAssignment) Get(db *gorm.DB, user *User) (int, *util.Response) {
	resp := new(util.Response) // Response placeholder

	assignment := new(Assignment) // Get the assignment by ID
	err := db.Where("id = ?", g.ID).First(assignment).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "not found"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}
	db.Model(assignment).Association("Files").Find(&assignment.Files) // Get all the files for the assignment

	// Get the subject the assignment is in
	subject := new(Subject)
	err = db.Where("id = ?", assignment.SubjectID).First(subject).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "not found"
			return 400, resp
		}
		return util.DatabaseError(err, resp)
	}
	// Get all students from the subject
	db.Model(subject).Association("Students").Find(&subject.Students)

	// If the user is a teacher then
	if user.Has(Teacher) {
		// Get all the requests, for each request get the requests,
		// and for each request get all the files that students uploaded
		db.Model(assignment).Association("Requests").Find(&assignment.Requests)
		for _, req := range assignment.Requests {
			upl := make([]*RequestUpload, 0)
			req.Uploads = upl
			uplBuf := make([]*RequestUpload, 0)
			db.Model(req).Association("Uploads").Find(&req.Uploads)
			req.Complete = nil

			for _, uploads := range req.Uploads {
				u := new(User)
				db.Where("id = ?", uploads.UserID).First(u)
				uploads.User = u
				uplBuf = append(uplBuf, uploads)
			}

			req.Uploads = uplBuf
		}
		// Get who has completed each assignment
		db.Model(assignment).Association("CompletedBy").Find(&assignment.CompletedBy)

		// If the assignment has a time due
		if assignment.TimeDue != nil {
			// Check who completed the assignment and who didnt
			completed := make(map[string]bool)
			for _, compl := range assignment.CompletedBy {
				completed[compl.ID] = true
			}

			assignment.NotCompletedBy = make([]*User, 0)
			for _, student := range subject.Students {
				if _, ok := completed[student.ID]; !ok {
					assignment.NotCompletedBy = append(assignment.NotCompletedBy, student)
				}
			}
		}
	} else {
		// Get all the requests for a student
		db.Model(assignment).Association("Requests").Find(&assignment.Requests)
		// For each request get what needs to be uploaded
		for r, req := range assignment.Requests {
			db.Model(req).Association("Uploads").Find(&assignment.Requests[r].Uploads)
			found := false
			for _, upload := range assignment.Requests[r].Uploads {
				if upload.UserID == user.ID {
					found = true
					break
				}
			}
			req.Complete = &found
			req.Uploads = nil
		}
	}

	resp.Data = assignment
	resp.Error = ""

	return 200, resp
}
