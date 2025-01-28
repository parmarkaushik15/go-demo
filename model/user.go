package model

import (
	"gorm.io/gorm"
)

// User struct
type User struct {
	gorm.Model
	ID           string  `json:"id" gorm:"primaryKey"`
    FirstName    string  `json:"first_name"`
    LastName     string  `json:"last_name"`
    EmailAddress string  `json:"email_address" gorm:"type:bytea"`
    CreatedAt    string  `json:"created_at"`
    DeletedAt    *string `json:"deleted_at,omitempty"`
    MergedAt     *string `json:"merged_at,omitempty"`
    ParentUserID *string `json:"parent_user_id,omitempty"`
}

// Users struct
type Users struct {
	Users []User `json:"users"`
}
 