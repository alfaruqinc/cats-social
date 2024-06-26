package domain

import (
	"log"
	"net/http"
)

type MessageErr interface {
	Message() string
	Status() int
	Error() string
}

type ErrorData struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

func (e *ErrorData) Message() string {
	return e.ErrMessage
}

func (e *ErrorData) Status() int {
	return e.ErrStatus
}

func (e *ErrorData) Error() string {
	return e.ErrError
}

func NewUnauthorizedError(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusForbidden,
		ErrError:   "NOT_AUTHORIZED",
	}
}

func NewUnauthenticatedError(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "NOT_AUTHENTICATED",
	}
}

func NewNotFoundError(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "NOT_FOUND",
	}
}

func NewBadRequest(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "BAD_REQUEST",
	}
}

func BadRequest(message string) error {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "BAD_REQUEST",
	}
}

func NewInternalServerError(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError, // 500
		ErrError:   "INTERNAL_SERVER_ERROR",
	}
}

func NewUnprocessibleEntityError(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "INVALID_REQUEST_BODY",
	}
}

func NewConflictError(message string) MessageErr {
	return &ErrorData{
		ErrMessage: message,
		ErrStatus:  http.StatusConflict,
		ErrError:   "CONFLICT_ERROR",
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}
}
