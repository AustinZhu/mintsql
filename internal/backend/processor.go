package backend

import "context"

type Processor interface {
	Process(ctx context.Context)
}
