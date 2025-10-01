package models

import "time"

type AuthProvider string

type User struct {
	ID           int64        `db:"id"`
	Email        string       `db:"email"`
	AuthProvider AuthProvider `db:"auth_provider"`
	PasswordHash *string      `db:"password_hash"` // nullable if auth provider is not "local"
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DisplayName  *string      `db:"display_name"` // nullable if user hasn't set it
}

type Waitlist struct {
	ID                 int64     `json:"id" db:"id"`
	Slug               string    `json:"slug" db:"slug"`
	Name               string    `json:"name" db:"name"`
	OwnerUserID        int64     `json:"owner_user_id" db:"owner_user_id"`
	IsPublic           bool      `json:"is_public" db:"is_public"`
	ShowVendorBranding bool      `json:"show_vendor_branding" db:"show_vendor_branding"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	ArchivedAt         *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}
