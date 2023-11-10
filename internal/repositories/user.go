package repositories

import (
	"context"
	"fmt"
	"strings"
	"user-management/internal/entities"
	"user-management/pkg/database"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create is an implementation of inserting a user entity
func (r *UserRepository) Create(ctx context.Context, db database.Executor, data *entities.User) error {
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

// GetUserByID is an implementation of retrieve user by id
func (r *UserRepository) GetUserByID(ctx context.Context, db database.Executor, id int64) (*entities.User, error) {
	var result entities.User
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

func (r *UserRepository) GetUserByUserName(ctx context.Context, db database.Executor, userName string) (*entities.User, error) {
	var result entities.User
	fieldNames, values := database.FieldMap(&result)
	stmt := fmt.Sprintf(`
		SELECT(%s) 
		FROM %s 
		WHERE user_name = $1
	`, strings.Join(fieldNames, ", "), result.TableName())
	row := db.QueryRowContext(ctx, stmt, userName)
	if err := row.Err(); err != nil {
		return nil, err
	}

	if err := row.Scan(values...); err != nil {
		return nil, err
	}

	return &result, nil
}
