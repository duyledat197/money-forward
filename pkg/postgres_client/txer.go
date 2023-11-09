package postgres_client

import (
	"context"
	"database/sql"
)

// Txer represents for a engine serve transaction
type Txer interface {
	Transaction(context.Context, func(context.Context, *sql.Tx) error) error
}

// txer is an implementation of Txer with postgres client
type txer struct {
	client *PostgresClient
}

func NewTxer(client *PostgresClient) Txer {
	return &txer{
		client: client,
	}
}

// Transaction implements a passing function with parameter have pointer of sql.Tx.
// The transaction begin with serializable isolation and then call passing function and then commit or rollback.
func (t *txer) Transaction(ctx context.Context, fn func(ctx context.Context, db *sql.Tx) error) error {
	tx, err := t.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(ctx, tx); err != nil {
		return err
	}

	tx.Commit()

	return nil
}
