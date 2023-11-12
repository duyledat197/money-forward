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

// GetUserByID is an implementation of retrieving user by id from database.
func (r *UserRepository) GetUserByID(ctx context.Context, db database.Executor, id int64) (*entities.UserWithAccounts, error) {
	userE := entities.User{}
	accountE := entities.Account{}
	fieldNames, _ := database.FieldMap(&userE)
	stmt := fmt.Sprintf(`
		SELECT %[2]s.%[1]s, 
		ARRAY_AGG(%[3]s.id) 
		FILTER(WHERE %[3]s.id IS NOT NULL) 
		FROM %[2]s 
		LEFT JOIN %[3]s ON %[2]s.id = %[3]s.user_id
		WHERE %[2]s.id = $1
		GROUP BY %[2]s.id,%[2]s.name
	`, strings.Join(fieldNames, ", users."), userE.TableName(), accountE.TableName())

	var result entities.UserWithAccounts
	row := db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	_, values := database.FieldMap(&result.User)

	if err := row.Scan(append(values, &result.AccountIDs)...); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUserByUserName is an implementation of retrieving user by userName from database.
func (r *UserRepository) GetUserByUserName(ctx context.Context, db database.Executor, userName string) (*entities.User, error) {
	var result entities.User
	fieldNames, values := database.FieldMap(&result)
	stmt := fmt.Sprintf(`
		SELECT %s
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

// UpdateByID is an implementation of updating user by id from database.
func (r *UserRepository) UpdateByID(ctx context.Context, db database.Executor, id int64, data *entities.User) error {
	e := &entities.User{}
	stmt := fmt.Sprintf(`
		UPDATE %s
		SET 
			name = COALESCE($2, name),
			updated_at = NOW()
		WHERE id = $1
	`, e.TableName())

	result, err := db.ExecContext(ctx, stmt, id, data.Name)
	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

// DeleteByID is an implementation of deleting user by id from database.
func (r *UserRepository) DeleteByID(ctx context.Context, db database.Executor, id int64) error {
	e := &entities.User{}
	stmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, e.TableName())

	result, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}
