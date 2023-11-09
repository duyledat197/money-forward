package postgres_client

import (
	"context"
	"database/sql"
)

// PostgresClient is presentation for a custom client of postgres.
type PostgresClient struct {
	*sql.DB
	connectionString string
}

func NewPostgresClient(connString string) *PostgresClient {
	return &PostgresClient{
		connectionString: connString,
	}
}

// Connect implements connect to postgres.
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

// Close implements close connection to postgres.
func (c *PostgresClient) Close(ctx context.Context) error {
	return c.DB.Close()
}
