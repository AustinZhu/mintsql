package parser

import (
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/lexer"
	"mintsql/internal/sql/token"
	"os"
)

type Parser struct {
	Ast   *ast.Ast
	Lexer *lexer.Lexer
}

func New(src string) *Parser {
	return &Parser{
		Ast:   &ast.Ast{},
		Lexer: lexer.New(src),
	}
}

func NewFromFile(path string) *Parser {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &Parser{
		Ast:   &ast.Ast{},
		Lexer: lexer.New(string(src)),
	}
}

func (p *Parser) Parse() error {
	go p.Lexer.Run()
	tokens := token.NewStream()
	delimiter := token.NewSymbol(token.SEMICOLON)

	for tk := p.Lexer.NextToken(); tk != nil; tokens = token.NewStream() {
		for tk != nil {
			tk = p.Lexer.NextToken()
			tokens.Add(tk)
			if tk.Equals(delimiter) {
				break
			}
		}
		if err := parseStmt(p.Ast, tokens); err != nil {
			return err
		}
	}
	return nil
}

func parseStmt(stmts *ast.Ast, tokens *token.Stream) error {
	init := tokens.Peek()
	if init.Kind != token.KindKeyword {
		return token.Error(init, "not a keyword", token.SELECT, token.INSERT, token.CREATE)
	}

	if token.NewKeyword(token.SELECT).Equals(init) {
		s, err := parseSelectStmt(tokens)
		if err != nil {
			return err
		}
		stmts.Add(s)
	}
	//if token.NewKeyword(token.CREATE).Equals(init) {
	//	return parseCreateStmt(p, tokens)
	//}
	//if token.NewKeyword(token.INSERT).Equals(init) {
	//	return parseInsertStmt(p, tokens)
	//}
	return token.Error(init, "unrecognized keyword", token.SELECT, token.INSERT, token.CREATE)
}

func parseSelectStmt(tokens *token.Stream) (*ast.Stmt, error) {
	stmt := &ast.Stmt{}

	if tk := tokens.Next(); !token.NewKeyword(token.SELECT).Equals(tk) {
		return nil, token.Error(tk, "not a select statement", token.SELECT)
	}
	stmt.Kind = ast.KindSelect
	stmt.SelectStmt = &ast.SelectStmt{}

	if exprs, err := parseExprs(tokens); err != nil {
		return nil, err
	} else {
		stmt.SelectStmt.Items = exprs
	}

	if tk := tokens.Next(); !token.NewKeyword(token.FROM).Equals(tk) {
		return nil, token.Error(tk, "missing keyword", token.FROM)
	}

	if tk := tokens.Next(); !token.HasKind(tk, token.KindIdentifier) {
		return nil, token.Error(tk, "invalid table name", "<identifier>")
	} else {
		stmt.SelectStmt.From = tk.Value
	}

	if tk := tokens.Next(); !token.NewSymbol(token.SEMICOLON).Equals(tk) {
		return nil, token.Error(tk, "missing ending semicolon", token.FROM)
	}

	if tk := tokens.Next(); tk != nil && tk.Kind != token.KindEof {
		return nil, token.Error(tk, "excessive tokens", "<nil>")
	}

	return stmt, nil
}

func parseExprs(tokens *token.Stream) ([]*ast.Expr, error) {
	expr := make([]*ast.Expr, 0)

	if tk := tokens.Next(); !token.HasKind(tk, token.KindIdentifier) {
		return nil, token.Error(tk, "invalid column name", "<identifier>")
	} else {
		expr = append(expr, &ast.Expr{Kind: ast.KindLiteral, Body: tk.Value})
	}

	for tk := tokens.Peek(); token.NewSymbol(token.COMMA).Equals(tk); {
		tk = tokens.Next()
		if tk = tokens.Next(); !token.HasKind(tk, token.KindIdentifier) {
			return nil, token.Error(tk, "invalid column name", "<identifier>")
		} else {
			expr = append(expr, &ast.Expr{Kind: ast.KindLiteral, Body: tk.Value})
		}
	}

	return expr, nil
}
