package db

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Storage defines the postgres storage
type Storage struct {
	Pgx    *pgxpool.Pool
	GormDB *gorm.DB
}

// NewStorage creates and returns a new Pgx storage connection
func NewStorage(connectionString string) (*Storage, error) {
	pgxConn, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	dbConnection, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	return &Storage{
		Pgx:    pgxConn,
		GormDB: dbConnection,
	}, nil
}

// Ping is a wrapper for Pgx Ping
func (s *Storage) Ping(ctx context.Context) error {
	return s.Pgx.Ping(ctx)
}

// Close all connections to database
func (s *Storage) Close() error {
	s.Pgx.Close()
	return nil
}
