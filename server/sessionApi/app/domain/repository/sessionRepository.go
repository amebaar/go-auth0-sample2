package repository

import (
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
)

type SessionRepository interface {
	InitState(ctx echo.Context) (string, error)
	IsValidState(ctx echo.Context, state string) bool
	Save(ctx echo.Context, userContext model.UserContext) error
	Discard(ctx echo.Context) error
}
