package database

import "database/sql"

// NullString help to transform string to [database/sql.NullString]
func NullString(str string) sql.NullString {
	var result sql.NullString
	result.Scan(str)

	return result
}
