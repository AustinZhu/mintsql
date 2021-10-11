package backend

import "mintsql/internal/sql/parser"

type QueryProcessor struct {
	Parser *parser.Parser
}

func (qp *QueryProcessor) Process() {

}
