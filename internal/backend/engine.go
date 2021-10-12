package backend

import (
	"context"
)

type Engine struct {
	Lang  QueryProcessor
	Store StoreProcessor
}

func Setup() *Engine {
	panic("TODO")
}

func (e *Engine) Execute(ctx context.Context, raw string) {

}
