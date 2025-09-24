package db

import (
	"context"

	"github.com/anish-chanda/openwaitlist/backend/internal/models"
)

type Database interface {

	// AUTH Stuff
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error

	// other helper functions
	Connect(dsn string) error
	Ping(ctx context.Context) error
	Close() error
	Migrate() error
}
