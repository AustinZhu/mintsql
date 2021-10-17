package parser

import (
	"github.com/google/go-cmp/cmp"
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/token"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input       string
		isError     bool
		expectedAst ast.Ast
	}{
		{
			input:   "CREATE TABLE users (id INT, name TEXT);\nINSERT INTO users VALUES (2, 'Kate');\nSELECT name, id FROM users;",
			isError: false,
			expectedAst: ast.Ast{
				{
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
				{
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
				{
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
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			parser := new(Parser)
			res, err := parser.Parse(test.input)
			if err != nil != test.isError {
				t.Error(err)
				return
			}
			if diff := cmp.Diff(test.expectedAst, res); diff != "" {
				t.Errorf("Parse mismatch (-want +got):\n%s", diff)
				return
			}
			return
		})
	}
}
