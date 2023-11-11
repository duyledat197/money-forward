package entities

import "database/sql"

type Account struct {
	ID        int64          `json:"id" db:"id"`
	Name      sql.NullString `json:"name" db:"name"`
	UserID    int64          `json:"user_id" db:"user_id"`
	Balance   sql.NullInt64  `json:"balance" db:"balance"`
	CreatedAt sql.NullTime   `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at" db:"updated_at"`
}

func (u *Account) TableName() string {
	return "accounts"
}
