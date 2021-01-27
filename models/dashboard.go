package models

// Dashboard is the user's dashboard
type Dashboard struct {
	User       *User                `json:"user"`
	Headmaster *HeadmasterDashboard `json:"headmaster,omitempty"`
	Teacher    *TeacherDashboard    `json:"teacher,omitempty"`
	Student    *StudentDashboard    `json:"student,omitempty"`
}

// HeadmasterDashboard contains headmaster information
type HeadmasterDashboard struct {
	Schools []*School `json:"schools"`
}

// TeacherDashboard contains headmaster information
type TeacherDashboard struct {
	SchoolID string     `json:"school_id"`
	Subjects []*Subject `json:"subjects"`
}

// StudentDashboard contains headmaster information
type StudentDashboard struct {
	Subject []*Subject `json:"subjects"`
}
