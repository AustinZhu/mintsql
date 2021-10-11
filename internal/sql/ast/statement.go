package ast

type SelectStmt struct {
	Table string
	Items []*Expr
}

type InsertStmt struct {
	Table  string
	Values []*Expr
}

type ColumnDef struct {
	Name     string
	DataType string
}

type CreateTableStmt struct {
	Name string
	Cols []*ColumnDef
}
