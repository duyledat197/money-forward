package deliveries

type CreateUserRequest struct {
	UserName string
	Password string
}
type CreateUserResponse struct {
	ID string
}

type GetUserByIDRequest struct{}
type GetUserByIDResponse struct{}
