package deliveries

import "context"

// using skeleton with cmd (d *userDelivery UserDelivery)
type UserDelivery interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUserByID(context.Context, *GetUserByIDRequest) (*GetUserByIDResponse, error)
}

type userDelivery struct {
}

func NewUserDelivery() UserDelivery {
	return &userDelivery{}
}

func (d *userDelivery) CreateUser(_ context.Context, _ *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{}, nil
}

func (d *userDelivery) GetUserByID(_ context.Context, _ *GetUserByIDRequest) (*GetUserByIDResponse, error) {
	return &GetUserByIDResponse{}, nil
}
