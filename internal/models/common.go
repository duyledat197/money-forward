package models

type Paging struct {
	Offset int64
	Limit  int64

	OrderBy   string
	OrderType string
}

type User struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	AccountIDs []int64 `json:"account_ids"`
}
