package backend

import (
	"context"
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/parser"
)

type QueryProcessor struct {
	Parser *parser.Parser
}

func (qp *QueryProcessor) Process(ctx context.Context, s string) (ast.Ast, error) {
	defer func() {
		qp.Parser = new(parser.Parser)
	}()
	return qp.Parser.Parse(s)
}
