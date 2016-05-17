package store

// DDL used to idempotently create db schema
const DDL = `
CREATE TABLE IF NOT EXISTS key_value (
	key TEXT PRIMARY KEY,
	value TEXT
);
`

func (s *Store) setValueForKey(key, value string) error {
	txn, err := s.querier.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare("INSERT OR REPLACE INTO key_value (key, value) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(key, value)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func (s *Store) valueFromKey(key string) (string, error) {
	var value string
	err := s.querier.
		QueryRow("SELECT value FROM key_value WHERE key = ?", key).
		Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}
