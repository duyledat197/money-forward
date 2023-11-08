package processor

import "context"

type Processor interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Factory interface {
	Connect(context.Context) error
	Close(context.Context) error
}
