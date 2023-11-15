package entities

import (
	"database/sql"

	"github.com/lib/pq"
)

type User struct {
	ID        int64          `json:"id" db:"id"`
	Name      sql.NullString `json:"name" db:"name"`
	UserName  string         `json:"user_name" db:"user_name"`
	Password  string         `json:"password" db:"password"`
	Role      User_Role      `json:"role" db:"role"`
	CreatedBy int64          `json:"created_by" db:"created_by"`
}

func (u *User) TableName() string {
	return "users"
}

// UserRole is the representation of a user role enum
type User_Role string

const (
	SuperAdminRole User_Role = "SUPER_ADMIN"
	AdminRole      User_Role = "ADMIN"
	UserRole       User_Role = "USER"
)

var (
	UserRoleList = []User_Role{SuperAdminRole, AdminRole, UserRole}
)

// UserWithAccounts is [User] extension with account id list inside.
type UserWithAccounts struct {
	User
	AccountIDs pq.Int64Array `json:"account_ids" db:"account_ids"`
}
