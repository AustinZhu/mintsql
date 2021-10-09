package ast

type ExprKind uint

const (
	KindLiteral ExprKind = iota
)

type Expr struct {
	Body string
	Kind ExprKind
}
