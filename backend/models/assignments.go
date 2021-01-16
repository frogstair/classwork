package models

import "time"

// Assignment is the internal model for assignments
type Assignment struct {
	ID        string    `json:"id"`
	TeacherID string    `json:"teacher_id"`
	Teacher   *User     `json:"teacher"`
	SubjectID string    `json:"subject_id"`
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	FilePath  string    `json:"file"`
	TimeDue   time.Time `json:"time_due"`
}
