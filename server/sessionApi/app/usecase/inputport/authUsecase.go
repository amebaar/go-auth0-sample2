package inputport

import (
	"github.com/labstack/echo/v4"
)

type AuthUsecase interface {
	InitState(ctx echo.Context) (string, error)
	Login(ctx echo.Context, request *AuthLoginRequest) error
}

type AuthLoginRequest struct {
	Username string
	Password string
	State    string
}
