package parser

import (
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/token"
)

func parseStmt(tokens *token.Stream) (*ast.Stmt, error) {
	init := tokens.Peek()
	if init.Kind != token.KindKeyword {
		return nil, Error(init, "not a keyword", token.SELECT, token.INSERT, token.CREATE)
	}

	var s *ast.Stmt
	var err error

	if token.NewKeyword(token.SELECT).Equals(init) {
		s, err = parseSelectStmt(tokens)
	} else if token.NewKeyword(token.CREATE).Equals(init) {
		s, err = parseCreateStmt(tokens)
	} else if token.NewKeyword(token.INSERT).Equals(init) {
		s, err = parseInsertStmt(tokens)
	} else {
		return nil, Error(init, "unrecognized keyword", token.SELECT, token.INSERT, token.CREATE)
	}

	if err != nil {
		return nil, err
	}
	return s, nil
}

func parseSelectStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := new(ast.Stmt)

	if tk := tokens.Next(); !token.NewKeyword(token.SELECT).Equals(tk) {
		return nil, Error(tk, "not a select statement", token.SELECT)
	}

	stmt.Kind = ast.KindSelect
	stmt.SelectStmt = &ast.SelectStmt{}

	if exprs, err := parseColumnExprs(tokens); err != nil {
		return nil, err
	} else {
		stmt.SelectStmt.Items = exprs
	}

	if tk := tokens.Next(); token.NewKeyword(token.FROM).NotEquals(tk) {
		return nil, Error(tk, "missing keyword", token.FROM)
	}

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, Error(tk, "invalid table name", "<identifier>")
	} else {
		stmt.SelectStmt.Table = tk.Value
	}

	if tk := tokens.Next(); token.NewSymbol(token.SEMICOLON).NotEquals(tk) {
		return nil, Error(tk, "missing ending semicolon", token.SEMICOLON)
	}

	if tk := tokens.Next(); token.NotEnd(tk) {
		return nil, Error(tk, "excessive tokens", "<nil>")
	}

	return stmt, nil
}

func parseCreateStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := new(ast.Stmt)

	if tk := tokens.Next(); token.NewKeyword(token.CREATE).NotEquals(tk) {
		return nil, Error(tk, "not a create statement", token.CREATE)
	}
	if tk := tokens.Next(); token.NewKeyword(token.TABLE).NotEquals(tk) {
		return nil, Error(tk, "missing keyword", token.TABLE)
	}

	stmt.Kind = ast.KindCreateTable
	stmt.CreateTableStmt = &ast.CreateTableStmt{}

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, Error(tk, "invalid table name", "<identifier>")
	} else {
		stmt.CreateTableStmt.Name = tk.Value
	}

	if tk := tokens.Next(); token.NewSymbol(token.LPAREN).NotEquals(tk) {
		return nil, Error(tk, "missing parentheses", token.LPAREN)
	}
	if defs, err := parseColumnDefs(tokens); err != nil {
		return nil, err
	} else {
		stmt.CreateTableStmt.Cols = defs
	}
	if tk := tokens.Next(); token.NewSymbol(token.RPAREN).NotEquals(tk) {
		return nil, Error(tk, "missing parentheses", token.RPAREN)
	}

	if tk := tokens.Next(); token.NewSymbol(token.SEMICOLON).NotEquals(tk) {
		return nil, Error(tk, "missing ending semicolon", token.SEMICOLON)
	}
	if tk := tokens.Next(); token.NotEnd(tk) {
		return nil, Error(tk, "excessive tokens", "<nil>")
	}

	return stmt, nil
}

func parseInsertStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := new(ast.Stmt)

	if tk := tokens.Next(); token.NewKeyword(token.INSERT).NotEquals(tk) {
		return nil, Error(tk, "not an insert statement", token.INSERT)
	}
	if tk := tokens.Next(); token.NewKeyword(token.INTO).NotEquals(tk) {
		return nil, Error(tk, "missing keyword", token.INTO)
	}

	stmt.Kind = ast.KindInsert
	stmt.InsertStmt = &ast.InsertStmt{}

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, Error(tk, "invalid table name", "<identifier>")
	} else {
		stmt.InsertStmt.Table = tk.Value
	}

	if tk := tokens.Next(); token.NewKeyword(token.VALUES).NotEquals(tk) {
		return nil, Error(tk, "missing keyword", token.VALUES)
	}

	if tk := tokens.Next(); token.NewSymbol(token.LPAREN).NotEquals(tk) {
		return nil, Error(tk, "missing parentheses", token.LPAREN)
	}
	if exprs, err := parseLiteralExprs(tokens); err != nil {
		return nil, err
	} else {
		stmt.InsertStmt.Values = exprs
	}
	if tk := tokens.Next(); token.NewSymbol(token.RPAREN).NotEquals(tk) {
		return nil, Error(tk, "missing parentheses", token.RPAREN)
	}

	if tk := tokens.Next(); token.NewSymbol(token.SEMICOLON).NotEquals(tk) {
		return nil, Error(tk, "missing ending semicolon", token.SEMICOLON)
	}
	if tk := tokens.Next(); token.NotEnd(tk) {
		return nil, Error(tk, "excessive tokens", "<nil>")
	}

	return stmt, nil
}
