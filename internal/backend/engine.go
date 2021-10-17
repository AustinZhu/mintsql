package backend

import (
	"context"
	"log"
	"mintsql/internal/sql/parser"
	"mintsql/internal/store/table"
	"time"
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

func (e *Engine) Execute(ctx context.Context, raw string) (res *table.Result, err error) {
	log.Println("Starting execution for request:", ctx.Value("uuid"))
	start := time.Now()
	defer func() {
		if err != nil {
			log.Printf("(%s) Execution failed: %s\n", ctx.Value("uuid"), err.Error())
			return
		}
		log.Printf("(%s) Execution completed: %s elapsed\n", ctx.Value("uuid"), time.Since(start))
	}()

	ast, err := e.Lang.Process(ctx, raw)
	if err != nil {
		return nil, err
	}
	return e.Store.Process(ctx, ast)
}
