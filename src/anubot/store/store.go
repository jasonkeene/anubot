package store

import (
	"database/sql"
	"log"
	"os/user"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

//go:generate hel -t Querier -o mock_querier_test.go

// DDL used to idempotently create db schema
const DDL = `
CREATE TABLE IF NOT EXISTS key_value (
	key TEXT PRIMARY KEY,
	value TEXT
);
`

// Querier is what is required by store to do its business, basically a sql.DB
type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
	Close() error
}

// Store is the the primary way of storing and retrieving data
type Store struct {
	querier Querier
}

// New creates a Store from a file path
func New(path string) *Store {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Panic(err)
	}
	return NewFromQuerier(db)
}

// NewFromQuerier creates a Store from a Querier like a sql.DB instance
func NewFromQuerier(querier Querier) *Store {
	return &Store{querier: querier}
}

// InitDDL executes idempotent DDL on the querier to make sure the schema is
// setup correctly
func (s *Store) InitDDL() error {
	txn, err := s.querier.Begin()
	if err != nil {
		return err
	}
	_, err = txn.Exec(DDL)
	if err != nil {
		return err
	}
	return txn.Commit()
}

// Close tearsdown all store related stuff
func (s *Store) Close() error {
	return s.querier.Close()
}

// SetCredentials stores credentials
func (s *Store) SetCredentials(user, pass string) error {
	txn, err := s.querier.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare("INSERT INTO key_value (key, value) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec("username", user)
	if err != nil {
		return err
	}
	_, err = stmt.Exec("password", pass)
	if err != nil {
		return err
	}
	err = txn.Commit()
	return err
}

// Credentials retrieves credentials
func (s *Store) Credentials() (string, string, error) {
	var user, pass string
	err := s.querier.
		QueryRow("SELECT value FROM key_value WHERE key = 'username'").
		Scan(&user)
	if err != nil {
		return "", "", err
	}
	err = s.querier.
		QueryRow("SELECT value FROM key_value WHERE key = 'password'").
		Scan(&pass)
	if err != nil {
		return "", "", err
	}
	return user, pass, nil
}

// HomePath returns a path for a sqlite db in the current user's home directory
func HomePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return path.Join(usr.HomeDir, "anubot.db")
}
