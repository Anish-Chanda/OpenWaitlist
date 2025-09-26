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
	ID                 int64      `db:"id"`
	Slug               string     `db:"slug"` // This is a Varchar(6), our app would generate a 6 Char alphanumeric string
	Name               string     `db:"name"`
	OwnerUserID        int64      `db:"owner_user_id"` //This is the foreign key to User.ID
	IsPublic           bool       `db:"is_public"`
	ShowVendorBranding bool       `db:"show_vendor_branding"` //if set to false, we hide badges/mentions of open-waitlist
	CreatedAt          time.Time  `db:"created_at"`
	ArchivedAt         *time.Time `db:"archived_at"`
}
