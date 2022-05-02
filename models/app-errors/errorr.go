package apperrors

import (
	"errors"
	"fmt"
)

type Type string

const (
	SameObj Type = "SAME_OBJECT"
)

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

/*
 * Error Üretme
 */

func NewSameObject(reason string) *Error {
	return &Error{
		Type:    SameObj,
		Message: fmt.Sprintf("Same Object: %v", reason),
	}
}

var ErrSameObj = errors.New("herhangi bir güncelleme mevcut değil")
