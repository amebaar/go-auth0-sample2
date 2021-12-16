package inputport

import (
	"github.com/labstack/echo/v4"
)

type AuthUsecase interface {
	InitState(ctx echo.Context) (string, error)
	Session(ctx echo.Context) (string, error)
	SignUp(ctx echo.Context, request *SignUpRequest) error
	Login(ctx echo.Context, request *AuthLoginRequest) error
	SocialLogin(ctx echo.Context, request *SocialLoginRequest) (string, error)
	UpdateSession(ctx echo.Context, request *UpdateSessionRequest) error
	Logout(ctx echo.Context) error
}

type SignUpRequest struct {
	Email    string
	Password string
}

type AuthLoginRequest struct {
	Username string
	Password string
	State    string
}

type SocialLoginRequest struct {
	Connection string
	State      string
	RedirectTo string
}

type UpdateSessionRequest struct {
	Code  string
	State string
}
