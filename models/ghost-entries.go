package models

import (
	"fmt"
)

// GhostEntry for maintaining Ghost entries
type GhostEntry struct {
	Ghostkey   interface{}
	Ghostvalue interface{}
	ghost      bool
}

func (e *GhostEntry) setLRU(l interface{}) {
	// e.detach()
	// e.ll = l.(*list.List)
	// e.el = e.ll.PushBack(e)
}

func (e *GhostEntry) setMRU(l interface{}) {
	// e.detach()
	// e.ll = l.(*list.List)
	// e.el = e.ll.PushFront(e)
}

func (e *GhostEntry) detach() {
	// if e.ll != nil {
	// 	e.ll.Remove(e.el)
	// }
}

// Get returns a list of ghost entries from database for the given list (B1, or B2)
func (g *GhostEntry) Get(database *Database, listID int) (map[interface{}]interface{}, error) {
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
func (g *GhostEntry) Add(database *Database, listID int, key, value interface{}) error {

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
