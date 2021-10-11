package backend

import "context"

type Engine struct {
	QueryProcessor QueryProcessor
	//Store
}

func Setup() *Engine {
	panic("TODO")
}

func (e *Engine) Execute(ctx context.Context, raw string) {

}
