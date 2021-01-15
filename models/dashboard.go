package models

// Dashboard is the user's dashboard
type Dashboard struct {
	Headmaster HeadmasterDashboard `json:"headmaster,omitempty"`
	Teacher    TeacherDashboard    `json:"teacher,omitempty"`
	Student    StudentDashboard    `json:"student,omitempty"`
}

// HeadmasterDashboard contains headmaster information
type HeadmasterDashboard struct {
}

// TeacherDashboard contains headmaster information
type TeacherDashboard struct {
}

// StudentDashboard contains headmaster information
type StudentDashboard struct {
}
