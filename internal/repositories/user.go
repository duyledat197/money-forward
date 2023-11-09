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
