package db

import (
	"context"

	"github.com/anish-chanda/openwaitlist/backend/internal/models"
)

type Database interface {

	// AUTH Stuff
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error

	// WAITLIST Stuff
	GetWaitlistsByUserID(ctx context.Context, userID int64, searchName string) ([]*models.Waitlist, error)
	CreateWaitlist(ctx context.Context, waitlist *models.Waitlist) error
	GetWaitlistByID(ctx context.Context, id int64) (*models.Waitlist, error)
	GetWaitlistBySlug(ctx context.Context, slug string) (*models.Waitlist, error)
	UpdateWaitlist(ctx context.Context, waitlist *models.Waitlist) error
	DeleteWaitlist(ctx context.Context, id int64) error

	// other helper functions
	Connect(dsn string) error
	Ping(ctx context.Context) error
	Close() error
	Migrate() error
}
