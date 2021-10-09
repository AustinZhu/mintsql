package ast

type SelectStmt struct {
	Items []*Expr
	From  string
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
