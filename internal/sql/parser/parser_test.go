package parser

import (
	"github.com/google/go-cmp/cmp"
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/lexer"
	"mintsql/internal/sql/token"
	"testing"
)

func TestParseSelectStmt(t *testing.T) {
	tests := []struct {
		input       string
		isError     bool
		expectedAst *ast.Stmt
	}{
		{
			input:   "SELECT name, id FROM users;",
			isError: false,
			expectedAst: &ast.Stmt{
				Kind: ast.KindSelect,
				SelectStmt: &ast.SelectStmt{
					Items: []*ast.Expr{
						{
							Body: "name",
							Kind: 0,
						},
						{
							Body: "id",
							Kind: 0,
						},
					},
					From: "users",
				},
				InsertStmt:      nil,
				CreateTableStmt: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			parser := New(test.input)
			tokens := extractTokens(parser.Lexer)
			selectStmt, err := parseSelectStmt(tokens)
			if err != nil != test.isError {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(test.expectedAst, selectStmt); diff != "" {
				t.Errorf("SelectStmt mismatch (-want +got):\n%s", diff)
				return
			}
			return
		})
	}
}

func extractTokens(l *lexer.Lexer) *token.Stream {
	stream := token.NewStream()
	go l.Run()
	for tk := l.NextToken(); tk != nil; tk = l.NextToken() {
		stream.Add(tk)
	}
	return stream
}
