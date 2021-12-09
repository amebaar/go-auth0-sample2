package middleware

import (
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/sdk/authman"
)

func HasPermission(operation, resource string) echo.MiddlewareFunc {
	a, _ := authman.GetAuthMan()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, err := a.LoadProfile(c); err != nil {
				c.Logger().Errorf("load profile error: %+v", err)
				c.Error(echo.ErrUnauthorized)
			}

			if a.HasPermission(c, operation, resource) {
				return next(c)
			} else {
				c.Error(echo.ErrForbidden)
			}
			return nil
		}
	}
}
