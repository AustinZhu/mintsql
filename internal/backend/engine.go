package backend

import (
	"context"
	"mintsql/internal/sql/parser"
	"mintsql/internal/store/table"
)

type Engine struct {
	Lang  *QueryProcessor
	Store *StoreProcessor
}

func Setup() *Engine {
	lang := &QueryProcessor{
		Parser: new(parser.Parser),
	}
	store := new(StoreProcessor)

	err := store.Load()
	if err != nil {
		panic(err)
	}
	return &Engine{
		Lang:  lang,
		Store: store,
	}
}

func (e *Engine) Execute(ctx context.Context, raw string) (*table.Result, error) {
	res, err := e.Lang.Process(ctx, raw)
	if err != nil {
		return nil, err
	}
	return e.Store.Process(ctx, res)
}
