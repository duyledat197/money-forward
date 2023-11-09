package entities

type User struct {
}

func (u *User) TableName() string {
	return "users"
}
