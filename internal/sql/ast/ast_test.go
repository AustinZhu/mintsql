package ast

import "testing"

func TestAst(t *testing.T) {
	stmt := new(Stmt)
	ast := New()
	ast.Add(stmt)
	if len(ast) != 1 {
		t.Errorf("ast.Len() = %d, want 1", len(ast))
	}
}
