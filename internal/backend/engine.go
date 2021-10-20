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

type ExecResult struct {
	err    error
	result *table.Result
}

func (er ExecResult) String() string {
	if er.err != nil {
		return er.err.Error()
	} else if er.result == nil {
		return "ok"
	} else {
		return er.result.String()
	}
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

func (e *Engine) Execute(ctx context.Context, raw string) (res *ExecResult) {
	log.Println("Starting execution for request:", ctx.Value("uuid"))
	start := time.Now()
	defer func() {
		if res.err != nil {
			log.Printf("(%s) Execution failed: %s\n", ctx.Value("uuid"), res.err.Error())
			return
		}
		log.Printf("(%s) Execution completed: %s elapsed\n", ctx.Value("uuid"), time.Since(start))
	}()

	ast, err := e.Lang.Process(ctx, raw)
	if err != nil {
		return &ExecResult{
			err:    err,
			result: nil,
		}
	}
	r, err := e.Store.Process(ctx, ast)
	return &ExecResult{
		err:    err,
		result: r,
	}
}
