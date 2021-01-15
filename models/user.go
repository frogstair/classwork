package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

// User is the internal representation of a user
type User struct {
	ID          string `gorm:"primaryKey"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Email       string `gorm:"not null;unique"`
	Password    string
	Token       string `gorm:"unique"`
	Perms       Role   `gorm:"not null"`
	PassSet     bool
	OneTimeCode string `gorm:"unique"`
}

// Has returns if a user has a role
func (u *User) Has(r Role) bool {
	role := u.Perms & r
	return role == 1
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
