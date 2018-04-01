package models

import (
	"database/sql"

	//"github.com/deepak11627/arc/arc"
	_ "github.com/go-sql-driver/mysql"
)

// Database to handle all DB operations
type Database struct {
	db *sql.DB
	//	logger arc.Logger
}

// SetLogger the logger
// func SetLogger(l arc.Logger) Option {
// 	return func(s *Database) error {
// 		s.logger = l
// 		return nil
// 	}
// }

// Option for optional params that can be set later
type Option func(s *Database) error

// NewDatabase returns a ShopStore instance
func NewDatabase(db *sql.DB, opts ...Option) *Database {

	s := &Database{db: db}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Open opens a new database.
func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn+"?parseTime=true")

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close closes the database.
func (s *Database) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
