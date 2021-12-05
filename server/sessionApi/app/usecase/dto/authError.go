package dto

import (
	"fmt"
)

type AuthBadRequestError struct {
	BaseErr error
}

func (e *AuthBadRequestError) Error() string {
	return fmt.Sprintf("Invalid Request: %+v", e.BaseErr)
}
func (e *AuthBadRequestError) Unwrap() error { return e.BaseErr }

type AuthUnauthorizedError struct {
	BaseErr error
}

func (e *AuthUnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %+v", e.BaseErr)
}
func (e *AuthUnauthorizedError) Unwrap() error { return e.BaseErr }

type AuthForbiddenError struct {
	BaseErr error
}

func (e *AuthForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden: %+v", e.BaseErr)
}
func (e *AuthForbiddenError) Unwrap() error { return e.BaseErr }

type AuthInternalServerError struct {
	BaseErr error
}

func (e *AuthInternalServerError) Error() string {
	return fmt.Sprintf("Internal Server Error: %+v", e.BaseErr)
}
func (e *AuthInternalServerError) Unwrap() error { return e.BaseErr }
