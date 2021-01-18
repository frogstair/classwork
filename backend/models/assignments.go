package models

import (
	"classwork/backend/util"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
)

// Assignment is the internal model for assignments
type Assignment struct {
	ID           string            `gorm:"primaryKey" json:"id"`
	TeacherID    string            `gorm:"not null" json:"teacher_id"`
	Teacher      *User             `gorm:"not null" json:"teacher,omitempty"`
	SubjectID    string            `gorm:"not null" json:"subject_id,omitempty"`
	Name         string            `gorm:"not null" json:"name"`
	Text         string            `json:"text"`
	Files        []*AssignmentFile `json:"files,omitempty"`
	TimeDue      *time.Time        `json:"time_due,omitempty"`
	TimeAssigned *time.Time        `json:"time_assigned"`
	Requests     []*Request        `json:"uploads,omitempty"`
	CompletedBy  []*User           `gorm:"many2many:assignments_completed" json:"comleted_by,omitempty"`
}

// AssignmentFile is the internal structure of a file relating to an assignment
type AssignmentFile struct {
	AssignmentID string `json:"-"`
	Path         string `gorm:"primaryKey" json:"path"`
}

// Request is the model to request an upload for the students
type Request struct {
	ID           string           `gorm:"primaryKey" json:"id"`
	AssignmentID string           `json:"-"`
	Name         string           `gorm:"not null" json:"name"`
	Uploads      []*RequestUpload `json:"uploads,omitempty"`
}

// RequestUpload is the model to tracks who uploaded what
type RequestUpload struct {
	UserID    string `json:"user"`
	RequestID string
	Filepath  string `json:"file"`
}

// NewAssignment is a model to create a new assignment
type NewAssignment struct {
	Name           string     `json:"name"`
	Text           string     `json:"text"`
	SubjectID      string     `json:"subject_id"`
	TimeDue        *time.Time `json:"time_due"`
	Files          []string   `json:"files"`
	UploadRequest  bool       `json:"-"`
	UploadRequests []struct {
		Name string `json:"name"`
	} `json:"uploads"`
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
func (n *NewAssignment) Create(db *gorm.DB, user *User) (int, *Response) {
	resp := new(Response)
	n.clean()

	if valid, reason := n.validate(); !valid {
		resp.Data = nil
		resp.Error = reason
		return 400, resp
	}

	if n.UploadRequest && (n.TimeDue == nil || n.TimeDue.Before(time.Now())) {
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
		file = util.ToRelativeFPath(file)
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

	now := time.Now()

	assignment := new(Assignment)
	assignment.ID = ksuid.New().String()
	assignment.Name = n.Name
	assignment.Text = n.Text
	assignment.TeacherID = user.ID
	assignment.SubjectID = subject.ID
	assignment.TimeDue = n.TimeDue
	assignment.TimeAssigned = &now

	err = db.Save(assignment).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	requests := make([]*Request, len(n.UploadRequests))

	for i, uploadReq := range n.UploadRequests {
		request := new(Request)
		request.ID = ksuid.New().String()
		request.AssignmentID = assignment.ID
		request.Name = uploadReq.Name

		requests[i] = request
		db.Save(request)
	}

	files := make([]*AssignmentFile, len(names))
	for i, name := range names {
		file := new(AssignmentFile)
		file.AssignmentID = assignment.ID
		file.Path = name

		files[i] = file
		db.Save(file)
	}

	assgn := struct {
		ID        string            `json:"id"`
		Name      string            `json:"name"`
		TeacherID string            `json:"teacher_id"`
		SubjectID string            `json:"subject_id"`
		Files     []*AssignmentFile `json:"files"`
		Uploads   []*Request        `json:"uploads"`
	}{assignment.ID, assignment.Name, assignment.TeacherID, assignment.SubjectID, files, requests}

	resp.Data = assgn
	resp.Error = ""

	return 200, resp
}

// NewRequestComplete is a model to complete an upload request
type NewRequestComplete struct {
	RequestID string `json:"request_id"`
	Filepath  string `json:"filepath"`
}

func (n *NewRequestComplete) clean() {
	util.RemoveSpaces(&n.Filepath)
	util.RemoveSpaces(&n.RequestID)
}

// Complete completes the upload request
func (n *NewRequestComplete) Complete(db *gorm.DB, user *User) (int, *Response) {
	resp := new(Response)

	n.clean()

	request := new(Request)
	err := db.Where("id = ?", n.RequestID).First(request).Error
	if err != nil {
		if util.IsNotFoundErr(err) {
			resp.Data = nil
			resp.Error = "request not found"
			return 400, resp
		}
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	assgn := new(Assignment)
	err = db.Where("id = ?", request.AssignmentID).First(assgn).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		log.Printf("Database error: %s\n", err.Error())
		return 500, resp
	}

	db.Model(request).Association("Uploads").Find(&request.Uploads)

	if _, err := os.Stat(util.ToRelativeFPath(n.Filepath)); os.IsNotExist(err) {
		resp.Data = nil
		resp.Error = "Internal error"
		return 500, resp
	}

	reqUpl := new(RequestUpload)
	reqUpl.Filepath = n.Filepath
	reqUpl.RequestID = request.ID
	reqUpl.UserID = user.ID

	request.Uploads = append(request.Uploads, reqUpl)
	err = db.Save(request).Error
	if err != nil {
		resp.Data = nil
		resp.Error = "Internal error"
		return 500, resp
	}

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

	if completed {
		assgn.CompletedBy = append(assgn.CompletedBy, user)
		err = db.Save(assgn).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			return 500, resp
		}
	}

	resp.Data = true
	resp.Error = ""

	return 202, resp
}
