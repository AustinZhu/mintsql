package backend

import (
	"context"
	"mintsql/internal/sql/ast"
	"mintsql/internal/sql/parser"
)

type QueryProcessor struct {
	Parser *parser.Parser
	Ast    chan *ast.Ast
}

func (qp *QueryProcessor) Process(ctx context.Context) error {
	return nil
}
