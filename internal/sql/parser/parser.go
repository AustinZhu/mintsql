package parser

import (
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/lexer"
	"mintsql/internal/sql/token"
	"os"
)

type Parser struct {
	Lexer *lexer.Lexer
}

func NewFromFile(path string) *Parser {
	src, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &Parser{
		Lexer: lexer.New(string(src)),
	}
}

func (p *Parser) Parse(src string) (ast ast.Ast, err error) {
	p.Lexer = lexer.New(src)
	go p.Lexer.Lex()

	tokens := token.NewStream()
	delimiter := token.NewSymbol(token.SEMICOLON)

	for tk := p.Lexer.NextToken(); token.NotEnd(tk); tokens = token.NewStream() {
		for ; token.NotEnd(tk); tk = p.Lexer.NextToken() {
			tokens.Add(tk)
			if tk.Equals(delimiter) {
				break
			}
		}
		if err = parseStmt(ast, tokens); err != nil {
			return nil, err
		}
		tk = p.Lexer.NextToken()
	}
	return
}
