package models

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}
type CreateUserResponse struct {
	ID int64 `json:"id"`
}

type GetUserByIDRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type GetUserByIDResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	AccountIDs []int64 `json:"account_ids"`
}

type UpdateUserRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type UpdateUserResponse struct {
}

type CreateAccountByUserIDRequest struct {
	UserID  int64  `json:"user_id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}
type CreateAccountByUserIDResponse struct {
	ID int64 `json:"id"`
}

type ListAccountByUserIDRequest struct {
	UserID int64 `json:"user_id"`
}

type ListAccountByUserIDResponse []*Account
