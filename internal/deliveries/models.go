package deliveries

type CreateUserRequest struct {
	UserName string
	Password string
}
type CreateUserResponse struct {
	ID string
}

type GetUserByIDRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type GetUserByIDResponse struct{}
