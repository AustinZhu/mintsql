package ast

import "mintsql/internal/sql/token"

type SelectStmt struct {
	items []*Expr
	from  token.Token
}

type InsertStmt struct {
	table  token.Token
	values []*Expr
}

type CreateTableStmt struct {
}
