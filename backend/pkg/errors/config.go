package errors

import (
	"fmt"
)

type ConfigFileError struct {
	FilePath string
	Err      error
}

func (e *ConfigFileError) Error() string {
	return fmt.Sprintf("Error reading config file %s: %v", e.FilePath, e.Err)
}

type ConfigValueError struct {
	Key string
	Err error
}

func (e *ConfigValueError) Error() string {
	return fmt.Sprintf("Error with config value for key %s: %v", e.Key, e.Err)
}
