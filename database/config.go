package database

import (
	"database/sql"
)

// ConfigEntry represents a configuration key-value pair
type ConfigEntry struct {
	Key   string
	Value string
}

// InitConfigTable creates the config table if it doesn't exist
func (db *DB) InitConfigTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	`
	_, err := db.conn.Exec(query)
	return err
}

// GetConfig retrieves a configuration value by key
func (db *DB) GetConfig(key string) (string, error) {
	var value string
	err := db.conn.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Return empty string if config doesn't exist
		}
		return "", err
	}
	return value, nil
}

// SetConfig sets a configuration value
func (db *DB) SetConfig(key, value string) error {
	_, err := db.conn.Exec(
		"INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)",
		key, value,
	)
	return err
}

// DeleteConfig deletes a configuration entry
func (db *DB) DeleteConfig(key string) error {
	_, err := db.conn.Exec("DELETE FROM config WHERE key = ?", key)
	return err
}
