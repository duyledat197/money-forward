package models

type GetAccountByIDRequest struct {
	ID int64
}

type GetAccountByIDResponse struct {
	*Account
}
