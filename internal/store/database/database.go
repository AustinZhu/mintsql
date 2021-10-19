package database

import (
	"mintsql/internal/store/table"
)

type Database struct {
	Tables map[string]*table.Table
}

func (db *Database) addTable(name string, tb *table.Table) error {
	if _, exists := db.Tables[name]; exists {
		return Error(DuplicateTableError, name, name)
	}
	db.Tables[name] = tb
	return nil
}

func (db *Database) getTable(name string) (*table.Table, error) {
	if tb, exists := db.Tables[name]; exists {
		return tb, nil
	}
	return nil, Error(NoSuchTableError, name, name)
}

func (db *Database) removeTable(name string) error {
	if _, exists := db.Tables[name]; exists {
		delete(db.Tables, name)
		return nil
	}
	return Error(DuplicateTableError, name, name)
}

func New() *Database {
	return &Database{Tables: make(map[string]*table.Table)}
}
