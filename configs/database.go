package configs

import "fmt"

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Schema   string
}

func (e *Database) Address() string {
	return fmt.Sprintf("%s:%s", e.Host, e.Port)
}
