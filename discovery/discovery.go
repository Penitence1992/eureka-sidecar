package discovery

import fmtfmt "fmt"

type Discovery interface {
	IsAppExists() (bool, error)
	CreateInstance() (bool, error)
	Heartbeat() (bool, error)
	RemoveInstance() (bool, error)
}

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
	return fmtfmt.Sprintf("%s, 错误码: %d", e.text, e.Code)
}
