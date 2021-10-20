package database

import (
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/token"
	"mintsql/internal/store/table"
)

func (db *Database) CreateTable(stmt *ast.CreateTableStmt) error {
	cols := make([]*table.Column, len(stmt.Cols))

	for i, c := range stmt.Cols {
		cols[i] = table.NewColumn(c.Name, c.DataType)
	}

	tb := table.NewTable(cols)
	err := db.addTable(stmt.Name, tb)

	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Insert(stmt *ast.InsertStmt) error {
	tb, err := db.getTable(stmt.Table)
	if err != nil {
		return err
	}

	row := make([]table.Cell, len(stmt.Values))

	for i, v := range stmt.Values {
		if v.Kind != ast.KindLiteral {
			return Error(IllegalValueError, stmt.Table, v.Body.Raw)
		}
		b := v.Body
		switch b.Kind {
		case token.KindNumeric:
			row[i] = table.FromString(b.Raw, table.Int)
		case token.KindString:
			row[i] = table.FromString(b.Raw, table.Text)
		default:
			return Error(NoSuchTypeError, stmt.Table, b.Raw)
		}
	}

	tb.Rows = append(tb.Rows, row)
	return nil
}

func (db *Database) Selects(stmt *ast.SelectStmt) (*table.Result, error) {
	tb, err := db.getTable(stmt.Table)
	if err != nil {
		return nil, err
	}

	if len(stmt.Items) == 1 && stmt.Items[0].Body.Raw == string(token.ASTERISK) {
		results := &table.Result{
			Columns: make([]table.Column, len(tb.Columns)),
			Rows:    make([][]table.Cell, len(tb.Rows)),
		}
		for i, c := range tb.Columns {
			results.Columns[i] = *c
		}
		results.Rows = tb.Rows
		return results, nil
	}

	results := &table.Result{
		Columns: make([]table.Column, len(stmt.Items)),
		Rows:    make([][]table.Cell, len(tb.Rows)),
	}

	idx := make([]int, len(stmt.Items))
	for i, c := range stmt.Items {
		if c.Kind != ast.KindColumn {
			return nil, Error(NoSuchColumnError, stmt.Table, c.Body.Raw)
		}
		for j, col := range tb.Columns {
			if col.Name == c.Body.Raw {
				idx[i] = j
			}
		}
	}

	for i, j := range idx {
		results.Columns[i] = *tb.Columns[j]
	}

	for i, r := range tb.Rows {
		results.Rows[i] = make([]table.Cell, len(idx))
		for j, k := range idx {
			results.Rows[i][j] = r[k]
		}
	}
	return results, nil
}
