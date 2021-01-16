package models

import (
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
func (u *User) GetDashboard(db *gorm.DB) (int, *Response) {
	resp := new(Response)

	dashboard := new(Dashboard)

	if u.Has(Headmaster) {
		hmDashboard := new(HeadmasterDashboard)

		err := db.Where("user_id = ?", u.ID).Find(&hmDashboard.Schools).Error
		if err != nil {
			resp.Data = nil
			resp.Error = "Internal error"
			log.Printf("Database error: %s\n", err.Error())
			return 500, resp
		}

		dashboard.Headmaster = hmDashboard
	}

	resp.Data = dashboard
	resp.Error = ""

	return 200, resp
}

// Delete deletes a user
func (u *User) Delete(db *gorm.DB) (int, *Response) {
	resp := new(Response)

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
func (u *User) Logout(db *gorm.DB) (int, *Response) {
	resp := new(Response)

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
