package models

import "github.com/jinzhu/gorm"

// User is the internal representation of a user
type User struct {
	ID          string `gorm:"primaryKey"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Email       string `gorm:"not null;unique"`
	Password    string `gorm:"not null"`
	Token       string `gorm:"unique"`
	Perms       Role   `gorm:"not null"`
	PassSet     bool
	OneTimeCode string `gorm:"unique"`
}

// Has returns if a user has a role
func (u *User) Has(r Role) bool {
	return u.Perms&r == 1
}

// GetDashboard gets the users dashboard
func (u *User) GetDashboard(db *gorm.DB) (int, *Response) {
	resp := new(Response)

	dashboard := new(Dashboard)

	resp.Data = dashboard
	resp.Error = ""

	return 200, resp
}
