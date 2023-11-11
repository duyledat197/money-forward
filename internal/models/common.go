package models

type User struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	AccountIDs []int64 `json:"account_ids"`
}

type Account struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}
