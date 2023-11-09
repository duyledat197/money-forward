package deliveries

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
type CreateUserResponse struct {
	ID string `json:"id"`
}

type GetUserByIDRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type GetUserByIDResponse struct {
}
