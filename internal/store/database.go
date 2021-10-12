package store

import "fmt"

type Database struct {
	Tables map[string]*Table
}

func (db *Database) AddTable(name string, tb *Table) error {
	if _, exists := db.Tables[name]; exists {
		return fmt.Errorf("table '%s' already exists", name)
	}
	db.Tables[name] = tb
	return nil
}

func (db *Database) GetTable(name string) (*Table, error) {
	if tb, exists := db.Tables[name]; exists {
		return tb, nil
	}
	return nil, fmt.Errorf("table '%s' does not exist", name)
}

func (db *Database) RemoveTable(name string) error {
	if _, exists := db.Tables[name]; exists {
		delete(db.Tables, name)
		return nil
	}
	return fmt.Errorf("table '%s' does not exist", name)
}

func NewDatabase() *Database {
	return &Database{Tables: make(map[string]*Table)}
}
