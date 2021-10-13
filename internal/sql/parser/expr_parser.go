package parser

import (
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/token"
)

func parseColumnExprs(tokens *token.Stream) ([]*ast.Expr, error) {
	expr := make([]*ast.Expr, 0)

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, token.Error(tk, "invalid column name", "<identifier>")
	} else {
		expr = append(expr, &ast.Expr{Kind: ast.KindColumn, Body: &ast.ExprBody{
			Raw:  tk.Value,
			Kind: tk.Kind,
		}})
	}

	for tk := tokens.Peek(); token.NewSymbol(token.COMMA).Equals(tk); tk = tokens.Peek() {
		tk = tokens.Next()
		if tk = tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
			return nil, token.Error(tk, "invalid column name", "<identifier>")
		} else {
			expr = append(expr, &ast.Expr{Kind: ast.KindColumn, Body: &ast.ExprBody{
				Raw:  tk.Value,
				Kind: tk.Kind,
			}})
		}
	}

	return expr, nil
}

func parseLiteralExprs(tokens *token.Stream) ([]*ast.Expr, error) {
	expr := make([]*ast.Expr, 0)

	if tk := tokens.Next(); token.NotKind(tk, token.KindString, token.KindNumeric) {
		return nil, token.Error(tk, "not values", "<string|numeric>")
	} else {
		expr = append(expr, &ast.Expr{Kind: ast.KindLiteral, Body: &ast.ExprBody{
			Raw:  tk.Value,
			Kind: tk.Kind,
		}})
	}

	for tk := tokens.Peek(); token.NewSymbol(token.COMMA).Equals(tk); tk = tokens.Peek() {
		tk = tokens.Next()
		if tk = tokens.Next(); token.NotKind(tk, token.KindString, token.KindNumeric) {
			return nil, token.Error(tk, "not values", "<string|numeric>")
		} else {
			expr = append(expr, &ast.Expr{Kind: ast.KindLiteral, Body: &ast.ExprBody{
				Raw:  tk.Value,
				Kind: tk.Kind,
			}})
		}
	}

	return expr, nil
}

func parseColumnDefs(tokens *token.Stream) ([]*ast.ColumnDef, error) {
	expr := make([]*ast.ColumnDef, 0)

	col := tokens.Next()
	if token.NotKind(col, token.KindIdentifier) {
		return nil, token.Error(col, "not a column name", "<identifier>")
	}
	dt := tokens.Next()
	if token.NewKeyword(token.INT).NotEquals(dt) && token.NewKeyword(token.TEXT).NotEquals(dt) {
		return nil, token.Error(dt, "not a datatype", "<datatype>")
	}
	expr = append(expr, &ast.ColumnDef{Name: col.Value, DataType: dt.Value})

	for tk := tokens.Peek(); token.NewSymbol(token.COMMA).Equals(tk); tk = tokens.Peek() {
		tk = tokens.Next()
		col := tokens.Next()
		if token.NotKind(col, token.KindIdentifier) {
			return nil, token.Error(col, "not a column name", "<identifier>")
		}
		dt := tokens.Next()
		if token.NewKeyword(token.INT).NotEquals(dt) && token.NewKeyword(token.TEXT).NotEquals(dt) {
			return nil, token.Error(dt, "not a datatype", "<datatype>")
		}
		expr = append(expr, &ast.ColumnDef{Name: col.Value, DataType: dt.Value})
	}

	return expr, nil
}
