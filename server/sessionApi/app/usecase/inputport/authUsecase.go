package inputport

import (
	"github.com/labstack/echo/v4"
)

type AuthUsecase interface {
	InitState(ctx echo.Context) (string, error)
	Session(ctx echo.Context) (string, error)
	Login(ctx echo.Context, request *AuthLoginRequest) error
	Logout(ctx echo.Context) error
}

type AuthLoginRequest struct {
	Username string
	Password string
	State    string
}
