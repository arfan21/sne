package constant

import (
	"fmt"
	"net/http"
)

var (
	ErrURLNotFound    = ErrWithCode{HTTPCode: http.StatusBadRequest, Message: "URL not found"}
	ErrTooManyRequest = ErrWithCode{HTTPCode: http.StatusTooManyRequests, Message: "Too many request"}
)

type ErrWithCode struct {
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
}

func (e ErrWithCode) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.HTTPCode, e.Message)
}

type ErrsWithCode []ErrWithCode

func (e ErrsWithCode) Error() string {
	var messages []string
	var httpCode int
	for _, err := range e {
		httpCode = err.HTTPCode
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("HTTP %d: %s", httpCode, messages)
}
