package configs

import "fmt"

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (e *Database) Address() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		e.Host,
		e.Port,
		e.User,
		e.Password,
		e.Database,
	)
}
