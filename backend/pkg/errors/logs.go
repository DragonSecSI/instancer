package errors

import (
	"fmt"
)

type LogsFileError struct {
	FilePath string
	Err      error
}

func (e *LogsFileError) Error() string {
	return fmt.Sprintf("Error opening logs file %s: %v", e.FilePath, e.Err)
}

type LogsLevelError struct {
	Level string
	Err   error
}

func (e *LogsLevelError) Error() string {
	return fmt.Sprintf("Invalid log level %s", e.Level)
}
