package models

import (
	"fmt"
)

/*
type Element struct {
	next, prev *Element

	// The list to which this element belongs.
	list *GhostList

	// The value stored with this element.
	Value interface{}
}

// // Next returns the next list element or nil.
func (e *Element) Next() *Element {
	// if p := e.next; e.list != nil && p != &e.list.root {
	// 	return p
	// }
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	// if p := e.prev; e.list != nil && p != &e.list.root {
	// 	return p
	// }
	return nil
}*/

// GhostList for maintaining Ghost entries
type GhostList struct {
	ID string
	// root     Element
	// len      int
	database *Database
}

// NewGhostList return a new ghost list
func NewGhostList(db *Database) *GhostList {

	gl := &GhostList{
		database: db,
		//	len:      0,
	}
	gl.Reset()
	return gl
}

/*
// Back returns last element from database
func (gl *GhostList) Back() *Element {
	return &Element{}
}


// Front returns the first element from database
func (gl *GhostList) Front() *Element {
	return &Element{}
}

// Len return count of values from database
func (gl *GhostList) Len() int {
	return 0
}

// Remove delete an element from database
func (gl *GhostList) Remove(*Element) interface{} {
	return nil
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (gl *GhostList) PushFront(e interface{}) *Element {
	return nil
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (gl *GhostList) PushBack(v interface{}) *Element {
	return nil
}*/

// Get returns a list of ghost entries from database for the given list (B1, or B2)
func (gl *GhostList) Get(listID string) (map[interface{}]interface{}, error) {
	logger := gl.database.logger
	logger.Debug("Geting key value pair from database.")
	rows, err := gl.database.db.Query("SELECT ghost_key, ghost_value FROM ghost_entries = ?", listID)
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

// Add saves a ghost List into database
func (gl *GhostList) PushFront(listID string, key, value interface{}) error {
	logger := gl.database.logger
	logger.Debug("Saving a key value pair in database.")

	stmt, err := gl.database.db.Prepare("INSERT INTO `ghost_lists` (`list_id`, `ghost_key`, `ghost_value`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `ghost_value` = VALUES(`ghost_value`);")

	if err != nil {
		logger.Debug(fmt.Sprintf("Error preparing add ghost entry %s", err))
		return fmt.Errorf("Error preparing add ghost entry %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		listID,
		key,
		value,
	)

	if err != nil {
		logger.Debug(fmt.Sprintf("Error inserting ghost entries: %s", err))
		return fmt.Errorf("Error inserting ghost entries: %s", err)
	}
	return nil
}

func (gl *GhostList) Remove(listID string) error {
	logger := gl.database.logger
	logger.Debug("Deleting list from Database.")
	stmt, err := gl.database.db.Prepare("DELETE FROM `ghost_lists` WHERE `list_id` = ? Limit 1")

	if err != nil {
		logger.Debug(fmt.Sprintf("Error preparing delet ghost entries: %s", err))
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(listID)

	if err != nil {
		logger.Debug(fmt.Sprintf("Error deleting ghost entries: %s", err))
	}
	// if no err, then err will be nil
	return err
}

func (gl *GhostList) Reset() error {
	logger := gl.database.logger
	logger.Debug("Deleting lists from Database.")
	_, err := gl.database.db.Query("DELETE FROM `ghost_lists`;")

	if err != nil {
		logger.Debug(fmt.Sprintf("Error preparing delet ghost entries: %s", err))
		return err
	}
	logger.Debug("Database Reset.")
	// if no err, then err will be nil
	return nil
}
