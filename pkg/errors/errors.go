package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func New(code int, message string) AppError {
	return AppError{Code: code, Message: message}
}

func Wrap(code int, message string, err error) AppError {
	if err == nil {
		return New(code, message)
	}
	return New(code, fmt.Sprintf("%s: %v", message, err))
}
