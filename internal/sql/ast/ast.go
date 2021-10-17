package ast

type Kind uint

const (
	KindSelect Kind = iota
	KindInsert
	KindCreateTable
)

type Ast []*Stmt

func (a *Ast) Add(stmt *Stmt) {
	*a = append(*a, stmt)
}

func New() Ast {
	stmts := make(Ast, 0)
	return stmts
}

type Stmt struct {
	Kind            Kind
	SelectStmt      *SelectStmt
	InsertStmt      *InsertStmt
	CreateTableStmt *CreateTableStmt
}
