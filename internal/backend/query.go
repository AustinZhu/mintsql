package backend

import (
	"context"
	"mintsql/internal/sql/parser"
)

type QueryProcessor struct {
	Parser *parser.Parser
}

func (qp *QueryProcessor) Process(ctx context.Context) {

}
