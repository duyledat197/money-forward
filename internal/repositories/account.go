package repositories

import (
	"context"
	"fmt"
	"strings"
	"user-management/internal/entities"
	"user-management/pkg/database"
)

type AccountRepository struct {
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

// Create is an implementation of inserting a user entity
func (r *AccountRepository) Create(ctx context.Context, db database.Executor, data *entities.Account) error {
	fieldNames, values := database.FieldMap(data)
	placeHolder := database.GetPlaceholders(len(fieldNames))
	stmt := fmt.Sprintf(`
		INSERT INTO %s(%s) VALUES(%s)
	`, data.TableName(), strings.Join(fieldNames, ", "), placeHolder)

	if _, err := db.ExecContext(ctx, stmt, values...); err != nil {
		return err
	}

	return nil
}

// ListAccountByUserID is an implementation of listing user by id from database.
func (r *AccountRepository) ListAccountByUserID(ctx context.Context, db database.Executor, userID int64) ([]*entities.Account, error) {
	return nil, nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, db database.Executor, id int64) (*entities.Account, error) {
	var result entities.Account
	fieldNames, values := database.FieldMap(&result)
	stmt := fmt.Sprintf(`
		SELECT(%s) 
		FROM %s 
		WHERE id = $1
	`, strings.Join(fieldNames, ", "), result.TableName())
	row := db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	if err := row.Scan(values...); err != nil {
		return nil, err
	}

	return &result, nil
}
