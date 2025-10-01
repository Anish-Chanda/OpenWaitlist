package postgres

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/anish-chanda/openwaitlist/backend/internal/logger"
	"github.com/anish-chanda/openwaitlist/backend/internal/models"
	"github.com/anish-chanda/openwaitlist/backend/migrations"
	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	conn *pgx.Conn
	log  logger.ServiceLogger
}

func NewPostgresDB(log logger.ServiceLogger) *PostgresDB {
	return &PostgresDB{
		log: log,
	}
}

// auth functions
func (s *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection is not established")
	}

	query := `
		SELECT id, email, auth_provider, password_hash, created_at, updated_at, display_name
		FROM users 
		WHERE email = $1
	`
	
	row := s.conn.QueryRow(ctx, query, email)
	
	var user models.User
	var passwordHash *string
	var displayName *string
	
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.AuthProvider,
		&passwordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&displayName,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		s.log.Error("Error getting user by email: ", err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	
	user.PasswordHash = passwordHash
	user.DisplayName = displayName
	
	return &user, nil
}

func (s *PostgresDB) CreateUser(ctx context.Context, user *models.User) error {
	if s.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	query := `
		INSERT INTO users (email, auth_provider, password_hash, created_at, updated_at, display_name)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	
	err := s.conn.QueryRow(ctx, query,
		user.Email,
		user.AuthProvider,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
		user.DisplayName,
	).Scan(&user.ID)
	
	if err != nil {
		s.log.Error("Error creating user: ", err)
		return fmt.Errorf("error creating user: %w", err)
	}
	
	s.log.Debug(fmt.Sprintf("Created user with ID: %d", user.ID))
	return nil
}

// Waitlist functions
func (s *PostgresDB) GetWaitlistsByUserID(ctx context.Context, userID int64, searchName string) ([]*models.Waitlist, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection is not established")
	}

	query := `
		SELECT id, slug, name, owner_user_id, is_public, show_vendor_branding, created_at, archived_at
		FROM waitlists 
		WHERE owner_user_id = $1 AND archived_at IS NULL
	`
	args := []interface{}{userID}

	// Add search filter if provided
	if searchName != "" {
		query += " AND name ILIKE $2"
		args = append(args, "%"+searchName+"%")
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		s.log.Error("Error querying waitlists: ", err)
		return nil, fmt.Errorf("error querying waitlists: %w", err)
	}
	defer rows.Close()

	var waitlists []*models.Waitlist
	for rows.Next() {
		var waitlist models.Waitlist
		err := rows.Scan(
			&waitlist.ID,
			&waitlist.Slug,
			&waitlist.Name,
			&waitlist.OwnerUserID,
			&waitlist.IsPublic,
			&waitlist.ShowVendorBranding,
			&waitlist.CreatedAt,
			&waitlist.ArchivedAt,
		)
		if err != nil {
			s.log.Error("Error scanning waitlist row: ", err)
			return nil, fmt.Errorf("error scanning waitlist: %w", err)
		}
		waitlists = append(waitlists, &waitlist)
	}

	if err = rows.Err(); err != nil {
		s.log.Error("Error iterating waitlist rows: ", err)
		return nil, fmt.Errorf("error iterating waitlists: %w", err)
	}

	return waitlists, nil
}

func (s *PostgresDB) CreateWaitlist(ctx context.Context, waitlist *models.Waitlist) error {
	if s.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	query := `
		INSERT INTO waitlists (slug, name, owner_user_id, is_public, show_vendor_branding, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	waitlist.CreatedAt = time.Now()

	err := s.conn.QueryRow(ctx, query,
		waitlist.Slug,
		waitlist.Name,
		waitlist.OwnerUserID,
		waitlist.IsPublic,
		waitlist.ShowVendorBranding,
		waitlist.CreatedAt,
	).Scan(&waitlist.ID)

	if err != nil {
		s.log.Error("Error creating waitlist: ", err)
		return fmt.Errorf("error creating waitlist: %w", err)
	}

	s.log.Debug(fmt.Sprintf("Created waitlist with ID: %d", waitlist.ID))
	return nil
}

func (s *PostgresDB) GetWaitlistByID(ctx context.Context, id int64) (*models.Waitlist, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection is not established")
	}

	query := `
		SELECT id, slug, name, owner_user_id, is_public, show_vendor_branding, created_at, archived_at
		FROM waitlists 
		WHERE id = $1 AND archived_at IS NULL
	`

	row := s.conn.QueryRow(ctx, query, id)

	var waitlist models.Waitlist
	err := row.Scan(
		&waitlist.ID,
		&waitlist.Slug,
		&waitlist.Name,
		&waitlist.OwnerUserID,
		&waitlist.IsPublic,
		&waitlist.ShowVendorBranding,
		&waitlist.CreatedAt,
		&waitlist.ArchivedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("waitlist not found")
		}
		s.log.Error("Error getting waitlist by id: ", err)
		return nil, fmt.Errorf("error getting waitlist: %w", err)
	}

	return &waitlist, nil
}

func (s *PostgresDB) GetWaitlistBySlug(ctx context.Context, slug string) (*models.Waitlist, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("database connection is not established")
	}

	query := `
		SELECT id, slug, name, owner_user_id, is_public, show_vendor_branding, created_at, archived_at
		FROM waitlists 
		WHERE slug = $1 AND archived_at IS NULL
	`

	row := s.conn.QueryRow(ctx, query, slug)

	var waitlist models.Waitlist
	err := row.Scan(
		&waitlist.ID,
		&waitlist.Slug,
		&waitlist.Name,
		&waitlist.OwnerUserID,
		&waitlist.IsPublic,
		&waitlist.ShowVendorBranding,
		&waitlist.CreatedAt,
		&waitlist.ArchivedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("waitlist not found")
		}
		s.log.Error("Error getting waitlist by slug: ", err)
		return nil, fmt.Errorf("error getting waitlist: %w", err)
	}

	return &waitlist, nil
}

func (s *PostgresDB) UpdateWaitlist(ctx context.Context, waitlist *models.Waitlist) error {
	if s.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	query := `
		UPDATE waitlists 
		SET slug = $1, name = $2, is_public = $3, show_vendor_branding = $4
		WHERE id = $5 AND archived_at IS NULL
	`

	_, err := s.conn.Exec(ctx, query,
		waitlist.Slug,
		waitlist.Name,
		waitlist.IsPublic,
		waitlist.ShowVendorBranding,
		waitlist.ID,
	)

	if err != nil {
		s.log.Error("Error updating waitlist: ", err)
		return fmt.Errorf("error updating waitlist: %w", err)
	}

	s.log.Debug(fmt.Sprintf("Updated waitlist with ID: %d", waitlist.ID))
	return nil
}

func (s *PostgresDB) DeleteWaitlist(ctx context.Context, id int64) error {
	if s.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	// Soft delete by setting archived_at timestamp
	query := `
		UPDATE waitlists 
		SET archived_at = $1
		WHERE id = $2 AND archived_at IS NULL
	`

	_, err := s.conn.Exec(ctx, query, time.Now(), id)
	if err != nil {
		s.log.Error("Error deleting waitlist: ", err)
		return fmt.Errorf("error deleting waitlist: %w", err)
	}

	s.log.Debug(fmt.Sprintf("Deleted waitlist with ID: %d", id))
	return nil
}

// Helper functions
func (s *PostgresDB) Connect(dsn string) error {
	parsedDSN, err := pgx.ParseConfig(dsn)
	if err != nil {
		s.log.Error("Error parsing DSN: ", err)
		return err
	}
	s.log.Debug(fmt.Sprintf("Connecting to Postgres at: %s:%d", parsedDSN.Host, parsedDSN.Port))

	conn, err := pgx.Connect(context.TODO(), dsn)
	if err != nil {
		s.log.Error("Failed to connect to PostgreSQL", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	s.conn = conn
	s.log.Info("Successfully connected to PostgreSQL")
	return nil
}

func (s *PostgresDB) Ping(ctx context.Context) error {
	if s.conn == nil {
		return fmt.Errorf("database connection is not established")
	}
	if err := s.conn.Ping(ctx); err != nil {
		s.log.Error("Database ping failed: ", err)
		return fmt.Errorf("database ping failed: %w", err)
	}
	s.log.Debug("Database ping successful")
	return nil
}

func (s *PostgresDB) Migrate() error {
	if s.conn == nil {
		return fmt.Errorf("database connection is not established")
	}

	s.log.Info("Starting database migrations...")

	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	
	if _, err := s.conn.Exec(context.Background(), createMigrationsTable); err != nil {
		s.log.Error("Failed to create migrations table: ", err)
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get migration files
	migrationsFS, dirName, err := migrations.GetMigrationsFS("postgresql")
	if err != nil {
		s.log.Error("Failed to get migrations filesystem: ", err)
		return fmt.Errorf("failed to get migrations filesystem: %w", err)
	}

	// Read all migration files
	var migrationFiles []string
	err = fs.WalkDir(migrationsFS, dirName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".up.sql") {
			migrationFiles = append(migrationFiles, path)
		}
		return nil
	})
	if err != nil {
		s.log.Error("Failed to walk migration files: ", err)
		return fmt.Errorf("failed to walk migration files: %w", err)
	}

	// Sort migration files to ensure proper order
	sort.Strings(migrationFiles)

	// Get already applied migrations
	appliedMigrations := make(map[string]bool)
	rows, err := s.conn.Query(context.Background(), "SELECT version FROM schema_migrations")
	if err != nil {
		s.log.Error("Failed to query applied migrations: ", err)
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			s.log.Error("Failed to scan migration version: ", err)
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		appliedMigrations[version] = true
	}

	// Apply pending migrations
	for _, migrationFile := range migrationFiles {
		// Extract version from filename (e.g., "0001_users_and_waitlist_table.up.sql" -> "0001_users_and_waitlist_table")
		baseName := filepath.Base(migrationFile)
		version := strings.TrimSuffix(baseName, ".up.sql")

		if appliedMigrations[version] {
			s.log.Debug(fmt.Sprintf("Migration %s already applied, skipping...", version))
			continue
		}

		s.log.Info("Applying migration: " + version)

		// Read migration content
		content, err := fs.ReadFile(migrationsFS, migrationFile)
		if err != nil {
			s.log.Error("Failed to read migration file: ", err)
			return fmt.Errorf("failed to read migration file %s: %w", migrationFile, err)
		}

		// Execute migration in a transaction
		tx, err := s.conn.Begin(context.Background())
		if err != nil {
			s.log.Error("Failed to start transaction: ", err)
			return fmt.Errorf("failed to start transaction for migration %s: %w", version, err)
		}

		// Execute migration SQL
		if _, err := tx.Exec(context.Background(), string(content)); err != nil {
			tx.Rollback(context.Background())
			s.log.Error("Failed to execute migration: ", err)
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record migration as applied
		if _, err := tx.Exec(context.Background(), "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback(context.Background())
			s.log.Error("Failed to record migration: ", err)
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		// Commit transaction
		if err := tx.Commit(context.Background()); err != nil {
			s.log.Error("Failed to commit migration transaction: ", err)
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		s.log.Info("Successfully applied migration: " + version)
	}

	s.log.Info("All migrations completed successfully")
	return nil
}

func (s *PostgresDB) Close() error {
	if s.conn != nil {
		if err := s.conn.Close(context.TODO()); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		return nil
	}
	return nil
}
