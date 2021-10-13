package backend

import (
	"context"
	"errors"
	"fmt"
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/token"
	"mintsql/internal/store"
	"strings"
)

type StoreProcessor struct {
	db *store.Database
}

func (sp *StoreProcessor) Process(ctx context.Context, stmts ast.Ast) (*Result, error) {
	for _, s := range stmts {
		switch s.Kind {
		case ast.KindSelect:
			return sp.selects(s.SelectStmt)
		case ast.KindInsert:
			err := sp.insert(s.InsertStmt)
			return nil, err
		case ast.KindCreateTable:
			err := sp.createTable(s.CreateTableStmt)
			return nil, err
		default:
			return nil, errors.New("not implemented")
		}
	}
	return nil, nil
}

func (sp *StoreProcessor) Load() error {
	sp.db = store.NewDatabase()
	return nil
}

type Result struct {
	Columns []store.Column
	Rows    [][]store.Cell
}

func (r Result) String() string {
	sb := new(strings.Builder)
	for _, col := range r.Columns {
		sb.WriteString(fmt.Sprintf("| %s ", col.Name))
	}
	sb.WriteString("|")

	for i := 0; i < 20; i++ {
		sb.WriteString("=")
	}
	sb.WriteString("\n")

	for _, result := range r.Rows {
		sb.WriteString("|")
		for i, cell := range result {
			typ := r.Columns[i].Type
			s := ""
			switch typ {
			case store.Int:
				s = fmt.Sprintf("%d", cell.AsInt())
			case store.Text:
				s = cell.AsText()
			}
			sb.WriteString(fmt.Sprintf(" %s | ", s))
		}
		sb.WriteString("\n")
	}
	return sb.String()
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

func (sp *StoreProcessor) insert(stmt *ast.InsertStmt) error {
	tb, err := sp.db.GetTable(stmt.Table)
	if err != nil {
		return err
	}
	row := make([]store.Cell, len(stmt.Values))
	for i, v := range stmt.Values {
		if v.Kind != ast.KindLiteral {
			return fmt.Errorf("expect values")
		}
		b := v.Body
		switch b.Kind {
		case token.KindNumeric:
			row[i] = store.FromString(b.Raw, store.Int)
		case token.KindString:
			row[i] = store.FromString(b.Raw, store.Text)
		default:
			return fmt.Errorf("unclassified types")
		}
	}
	tb.Rows = append(tb.Rows, row)
	return nil
}

func (sp *StoreProcessor) selects(stmt *ast.SelectStmt) (*Result, error) {
	tb, err := sp.db.GetTable(stmt.Table)
	if err != nil {
		return nil, err
	}
	results := &Result{
		Columns: make([]store.Column, len(stmt.Items)),
		Rows:    make([][]store.Cell, 0),
	}
	idx := make([]int, len(stmt.Items))
	for i, c := range stmt.Items {
		if c.Kind != ast.KindColumn {
			return nil, fmt.Errorf("expect columns")
		}
		for j, col := range tb.Columns {
			if col.Name == c.Body.Raw {
				idx[i] = j
			}
		}
	}
	for i, j := range idx {
		results.Columns[i] = *tb.Columns[j]
		results.Rows[i] = tb.Rows[j]
	}
	return results, nil
}
