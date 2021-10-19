package parser

import (
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/token"
)

func parseStmt(tokens *token.Stream) (*ast.Stmt, error) {
	init := tokens.Peek()
	if init.Kind != token.KindKeyword {
		return nil, Error(InvalidStmtError, init, token.SELECT, token.INSERT, token.CREATE)
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
		return nil, Error(InvalidStmtError, init, token.SELECT, token.INSERT, token.CREATE)
	}

	if err != nil {
		return nil, err
	}
	return s, nil
}

func parseSelectStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := new(ast.Stmt)

	if tk := tokens.Next(); !token.NewKeyword(token.SELECT).Equals(tk) {
		return nil, Error(InvalidStmtError, tk, token.SELECT)
	}

	stmt.Kind = ast.KindSelect
	stmt.SelectStmt = &ast.SelectStmt{}

	if exprs, err := parseColumnExprs(tokens); err != nil {
		return nil, err
	} else {
		stmt.SelectStmt.Items = exprs
	}

	if tk := tokens.Next(); token.NewKeyword(token.FROM).NotEquals(tk) {
		return nil, Error(MissingKeywordError, tk, token.FROM)
	}

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, Error(InvalidNameError, tk)
	} else {
		stmt.SelectStmt.Table = tk.Value
	}

	if tk := tokens.Next(); token.NewSymbol(token.SEMICOLON).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.SEMICOLON)
	}

	if tk := tokens.Next(); token.NotEnd(tk) {
		return nil, Error(ExcessiveTokenError, tk)
	}

	return stmt, nil
}

func parseCreateStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := new(ast.Stmt)

	if tk := tokens.Next(); token.NewKeyword(token.CREATE).NotEquals(tk) {
		return nil, Error(InvalidStmtError, tk, token.CREATE)
	}
	if tk := tokens.Next(); token.NewKeyword(token.TABLE).NotEquals(tk) {
		return nil, Error(MissingKeywordError, tk, token.TABLE)
	}

	stmt.Kind = ast.KindCreateTable
	stmt.CreateTableStmt = &ast.CreateTableStmt{}

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, Error(InvalidNameError, tk)
	} else {
		stmt.CreateTableStmt.Name = tk.Value
	}

	if tk := tokens.Next(); token.NewSymbol(token.LPAREN).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.LPAREN)
	}
	if defs, err := parseColumnDefs(tokens); err != nil {
		return nil, err
	} else {
		stmt.CreateTableStmt.Cols = defs
	}
	if tk := tokens.Next(); token.NewSymbol(token.RPAREN).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.RPAREN)
	}

	if tk := tokens.Next(); token.NewSymbol(token.SEMICOLON).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.SEMICOLON)
	}
	if tk := tokens.Next(); token.NotEnd(tk) {
		return nil, Error(ExcessiveTokenError, tk)
	}

	return stmt, nil
}

func parseInsertStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := new(ast.Stmt)

	if tk := tokens.Next(); token.NewKeyword(token.INSERT).NotEquals(tk) {
		return nil, Error(InvalidStmtError, tk, token.INSERT)
	}
	if tk := tokens.Next(); token.NewKeyword(token.INTO).NotEquals(tk) {
		return nil, Error(MissingKeywordError, tk, token.INTO)
	}

	stmt.Kind = ast.KindInsert
	stmt.InsertStmt = &ast.InsertStmt{}

	if tk := tokens.Next(); token.NotKind(tk, token.KindIdentifier) {
		return nil, Error(InvalidNameError, tk)
	} else {
		stmt.InsertStmt.Table = tk.Value
	}

	if tk := tokens.Next(); token.NewKeyword(token.VALUES).NotEquals(tk) {
		return nil, Error(InvalidValueError, tk, token.VALUES)
	}

	if tk := tokens.Next(); token.NewSymbol(token.LPAREN).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.LPAREN)
	}
	if exprs, err := parseLiteralExprs(tokens); err != nil {
		return nil, err
	} else {
		stmt.InsertStmt.Values = exprs
	}
	if tk := tokens.Next(); token.NewSymbol(token.RPAREN).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.RPAREN)
	}

	if tk := tokens.Next(); token.NewSymbol(token.SEMICOLON).NotEquals(tk) {
		return nil, Error(MissingSymbolError, tk, token.SEMICOLON)
	}
	if tk := tokens.Next(); token.NotEnd(tk) {
		return nil, Error(ExcessiveTokenError, tk)
	}

	return stmt, nil
}
