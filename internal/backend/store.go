package backend

import (
	"context"
	"errors"
	"mintsql/internal/sql/ast"
	"mintsql/internal/store/database"
	"mintsql/internal/store/table"
)

type StoreProcessor struct {
	db *database.Database
}

func (sp *StoreProcessor) Process(ctx context.Context, stmts ast.Ast) (*table.Result, error) {
	for _, s := range stmts {
		switch s.Kind {
		case ast.KindSelect:
			return sp.db.Selects(s.SelectStmt)
		case ast.KindInsert:
			err := sp.db.Insert(s.InsertStmt)
			return nil, err
		case ast.KindCreateTable:
			err := sp.db.CreateTable(s.CreateTableStmt)
			return nil, err
		default:
			return nil, errors.New("not implemented")
		}
	}
	return nil, nil
}

func (sp *StoreProcessor) Load() error {
	sp.db = database.NewDatabase()
	return nil
}
