package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Client represents a PostgreSQL database client
type Client struct {
	db *sql.DB
}

// NewClient creates a new PostgreSQL client
func NewClient(cfg Config) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	return NewClientFromDSN(dsn)
}

// NewClientFromDSN creates a new PostgreSQL client from a connection string
func NewClientFromDSN(dsn string) (*Client, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &Client{db: db}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}

// DB returns the underlying *sql.DB instance
func (c *Client) DB() *sql.DB {
	return c.db
}
