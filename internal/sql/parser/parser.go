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
		Ast:   new(ast.Ast),
		Lexer: lexer.New(src),
	}
}

func NewFromFile(path string) *Parser {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &Parser{
		Ast:   new(ast.Ast),
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
