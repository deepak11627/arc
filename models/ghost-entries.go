package models

import (
	"fmt"
)

// DBService is a generic User operations interface
type DBService interface {
	Get(ListID int, key interface{}) (interface{}, error)
	Add(ListID int, key interface{}, value interface{}) error
}

// GhostEntries for maintaining Ghost entries
type GhostEntries struct {
	ListID     int
	Ghostkey   interface{}
	Ghostvalue interface{}
}

// Get returns a list of ghost entries from database for the given list (B1, or B2)
func (g *GhostEntries) Get(database *Database, listID int) (map[interface{}]interface{}, error) {
	rows, err := database.db.Query("SELECT ghost_key, ghost_value FROM ghost_entries = ?", listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make(map[interface{}]interface{}, 0)
	for rows.Next() {
		var key interface{}
		var value interface{}
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		entries[key] = value
	}

	return entries, nil
}

// Add saves a ghost entry into database
func (g *GhostEntries) Add(database *Database, listID int, key, value interface{}) error {

	stmt, err := database.db.Prepare("INSERT INTO `ghost_entries` (`list_id`, `ghost_key`, `ghost_value`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `ghost_value` = VALUES(`ghost_value`);")

	if err != nil {
		return fmt.Errorf("Error preparing add ghost entry %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		listID,
		key,
		value,
	)

	if err != nil {
		return fmt.Errorf("Error inserting ghost entries: %s", err)
	}
	return nil
}
