package errors

import (
	"fmt"
)

type DatabaseConnectionError struct {
	Dialect string
	DSN     string
	Err     error
}

func (e *DatabaseConnectionError) Error() string {
	return fmt.Sprintf("Error connecting to database with dialect '%s' and DSN '%s': %v", e.Dialect, e.DSN, e.Err)
}

type DatabaseQueryError struct {
	Query string
	Err   error
}

func (e *DatabaseQueryError) Error() string {
	return fmt.Sprintf("Error executing database query '%s': %v", e.Query, e.Err)
}
