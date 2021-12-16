package web

import (
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	authMiddle "go-auth0-sample2/sdk/authman/middleware"
	"net/http"

	"go-auth0-sample2/server/sessionApi/app/adapter/web/controller"
)

func newRouter(
	authController controller.AuthController,
) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	err := setMiddleware(e)
	if err != nil {
		panic(err)
	}

	err = setRoute(e, authController)
	if err != nil {
		panic(err)
	}

	return e
}

func setMiddleware(e *echo.Echo) error {
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Session Setup
	store, err := redistore.NewRediStore(10, "tcp", "redis:6379", "", []byte("secret"))
	if err != nil {
		return err
	}
	e.Use(session.Middleware(sessions.Store(store)))

	return nil
}

func setRoute(e *echo.Echo,
	authController controller.AuthController,
) error {
	e.GET("/state",
		func(c echo.Context) error { return authController.InitState(c) })
	e.GET("/session",
		func(c echo.Context) error { return authController.GetSession(c) })
	e.POST("/signup",
		func(c echo.Context) error { return authController.SignUp(c) })
	e.POST("/login",
		func(c echo.Context) error { return authController.Login(c) })
	e.GET("/login",
		func(c echo.Context) error { return authController.SocialLogin(c) })
	e.GET("/logout",
		func(c echo.Context) error { return authController.Logout(c) })
	e.GET("/callback",
		func(c echo.Context) error { return authController.Callback(c) })

	e.GET("/read",
		func(c echo.Context) error { return c.String(http.StatusOK, "readOk") },
		authMiddle.HasPermission("read", "sample"))
	e.GET("/create",
		func(c echo.Context) error { return c.String(http.StatusOK, "createOk") },
		authMiddle.HasPermission("create", "sample"))

	return nil
}
