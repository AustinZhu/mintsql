package database

import (
	"fmt"
	"mintsql/internal/store/table"
)

type Database struct {
	Tables map[string]*table.Table
}

func (db *Database) addTable(name string, tb *table.Table) error {
	if _, exists := db.Tables[name]; exists {
		return fmt.Errorf("table '%s' already exists", name)
	}
	db.Tables[name] = tb
	return nil
}

func (db *Database) getTable(name string) (*table.Table, error) {
	if tb, exists := db.Tables[name]; exists {
		return tb, nil
	}
	return nil, fmt.Errorf("table '%s' does not exist", name)
}

func (db *Database) removeTable(name string) error {
	if _, exists := db.Tables[name]; exists {
		delete(db.Tables, name)
		return nil
	}
	return fmt.Errorf("table '%s' does not exist", name)
}

func New() *Database {
	return &Database{Tables: make(map[string]*table.Table)}
}
