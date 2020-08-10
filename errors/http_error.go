package errors

import "fmt"

func NewError(text string, code int) error {
	return &HttpError{
		text: text,
		Code: code,
	}
}

type HttpError struct {
	text string
	Code int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s, 错误码: %d", e.text, e.Code)
}
