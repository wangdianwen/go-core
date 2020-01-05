package utils

import "fmt"

type BError struct {
	Error   error
	Message string
	Code    int
}

func (b BError) String() string {
	return fmt.Sprintf("Error: %v, Message: %s, Code: %d", b.Error, b.Message, b.Code)
}
