package ast

type Kind uint

const (
	KindSelect = iota
	KindInsert
	KindCreateTable
)

type Ast struct {
	Kind            Kind
	SelectStmt      *SelectStmt
	InsertStmt      *InsertStmt
	CreateTableStmt *CreateTableStmt
}
