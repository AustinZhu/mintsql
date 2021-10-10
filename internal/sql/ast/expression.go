package ast

type ExprKind uint

const (
	KindLiteral ExprKind = iota
	KindColumn
	KindAsterisk
	KindAggregate
)

type Expr struct {
	Body string
	Kind ExprKind
}
