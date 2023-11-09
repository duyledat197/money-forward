package postgres_client

import (
	"context"
	"database/sql"
)

type PostgresClient struct {
	*sql.DB
	connectionString string
}

func NewPostgresClient(connString string) *PostgresClient {
	return &PostgresClient{
		connectionString: connString,
	}
}

func (c *PostgresClient) Connect(ctx context.Context) error {
	var err error
	c.DB, err = sql.Open("postgres", c.connectionString)
	if err != nil {
		return err
	}

	if err := c.DB.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *PostgresClient) Close(ctx context.Context) error {
	return c.DB.Close()
}
