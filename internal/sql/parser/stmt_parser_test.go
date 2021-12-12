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
							Body: &ast.ExprBody{
								Raw:  "name",
								Kind: token.KindIdentifier,
							},
							Kind: ast.KindColumn,
						},
						{
							Body: &ast.ExprBody{
								Raw:  "id",
								Kind: token.KindIdentifier,
							},
							Kind: ast.KindColumn,
						},
					},
					Table: "users",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := lexer.New(test.input)
			tokens := extractTokens(l)
			selectStmt, err := parseSelectStmt(tokens)
			if err != nil != test.isError {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(test.expectedAst, selectStmt); diff != "" {
				t.Errorf("SelectStmt mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestParseInsertStmt(t *testing.T) {
	tests := []struct {
		input       string
		isError     bool
		expectedAst *ast.Stmt
	}{
		{
			input:   "INSERT INTO users VALUES (2, 'Kate');",
			isError: false,
			expectedAst: &ast.Stmt{
				Kind: ast.KindInsert,
				InsertStmt: &ast.InsertStmt{
					Table: "users",
					Values: []*ast.Expr{
						{
							Body: &ast.ExprBody{
								Raw:  "2",
								Kind: token.KindNumeric,
							},
							Kind: ast.KindLiteral,
						},
						{
							Body: &ast.ExprBody{
								Raw:  "Kate",
								Kind: token.KindString,
							},
							Kind: ast.KindLiteral,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := lexer.New(test.input)
			tokens := extractTokens(l)
			insertStmt, err := parseInsertStmt(tokens)
			if err != nil != test.isError {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(test.expectedAst, insertStmt); diff != "" {
				t.Errorf("InsertStmt mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestParseCreateStmt(t *testing.T) {
	tests := []struct {
		input       string
		isError     bool
		expectedAst *ast.Stmt
	}{
		{
			input:   "CREATE TABLE users (id INT, name TEXT);",
			isError: false,
			expectedAst: &ast.Stmt{
				Kind: ast.KindCreateTable,
				CreateTableStmt: &ast.CreateTableStmt{
					Name: "users",
					Cols: []*ast.ColumnDef{
						{
							Name:     "id",
							DataType: "INT",
						},
						{
							Name:     "name",
							DataType: "TEXT",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := lexer.New(test.input)
			tokens := extractTokens(l)
			createStmt, err := parseCreateStmt(tokens)
			if err != nil != test.isError {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(test.expectedAst, createStmt); diff != "" {
				t.Errorf("CreateStmt mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func extractTokens(l *lexer.Lexer) *token.Stream {
	stream := token.NewStream()
	go l.Lex()
	for tk := l.NextToken(); tk != nil; tk = l.NextToken() {
		stream.Add(tk)
	}
	return stream
}
