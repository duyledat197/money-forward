package database

import (
	"context"
	"database/sql"
)

// Executor is a presentation of an database executor with exec and query command.
// Example implementation is [database/sql.Tx] and [database/sql.DB]
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
