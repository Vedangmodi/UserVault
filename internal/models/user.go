package models

import (
	"time"
)

// User is the API representation of a user.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"` // ISO date string (YYYY-MM-DD)
	Age  int    `json:"age,omitempty"`
}

// CreateUserRequest is used for POST /users.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is used for PUT /users/:id.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// CalculateAge returns the age in whole years as of "now" for a given date of birth.
func CalculateAge(dob time.Time, now time.Time) int {
	year, month, day := dob.Date()
	nowYear, nowMonth, nowDay := now.Date()

	age := nowYear - year
	// If birthday hasn't occurred yet this year, subtract 1.
	if nowMonth < month || (nowMonth == month && nowDay < day) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}


