package backend

import (
	"context"
	"mintsql/internal/sql/ast"
	"mintsql/internal/store"
)

type StoreProcessor struct {
	db store.Database
}

func (sp *StoreProcessor) Process(ctx context.Context) error {
	return nil
}

func (sp *StoreProcessor) Load() error {
	return nil
}

func (sp *StoreProcessor) createTable(stmt *ast.CreateTableStmt) error {
	cols := make([]*store.Column, len(stmt.Cols))
	for i, c := range stmt.Cols {
		cols[i] = store.NewColumn(c.Name, c.DataType)
	}
	table := store.NewTable(cols)
	err := sp.db.AddTable(stmt.Name, table)
	if err != nil {
		return err
	}
	return nil
}
