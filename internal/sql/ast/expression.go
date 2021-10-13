package ast

import "mintsql/internal/sql/token"

type ExprKind uint

const (
	KindLiteral ExprKind = iota
	KindColumn
	KindAsterisk
	KindAggregate
)

type Expr struct {
	Body *ExprBody
	Kind ExprKind
}

type ExprBody struct {
	Raw  string
	Kind token.Kind
}
