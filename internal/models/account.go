package models

type Account struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type GetAccountByIDRequest struct {
	*Account
}

type GetAccountByIDResponse struct {
}
