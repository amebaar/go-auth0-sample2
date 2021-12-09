package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/server/sessionApi/app/usecase/dto"
	"go-auth0-sample2/server/sessionApi/app/usecase/inputport"
	"go-auth0-sample2/server/sessionApi/config"
	"net/http"
	"net/url"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type AuthController interface {
	InitState(ctx echo.Context) error
	GetSession(ctx echo.Context) error
	Login(ctx echo.Context) error
	Logout(ctx echo.Context) error
}

type authController struct {
	authUsecase inputport.AuthUsecase
}

func NewAuthController(authUsecase inputport.AuthUsecase) AuthController {
	validate = validator.New()
	return &authController{authUsecase}
}

type LoginRequest struct {
	Username   string `json:"username" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	State      string `json:"state" validate:"required"`
	RedirectTo string `json:"redirect_to" validate:"required,url"`
}

func (c *authController) GetSession(ctx echo.Context) error {
	profile, err := c.authUsecase.Session(ctx)
	if err != nil {
		ctx.Logger().Infof("Failed to get session: %+v", err)
		return transDtoErrorToEcho(err)
	}
	result := map[string]string{"profile": profile}
	return ctx.JSON(http.StatusOK, result)
}

func (c *authController) InitState(ctx echo.Context) error {
	state, err := c.authUsecase.InitState(ctx)
	if err != nil {
		ctx.Logger().Infof("Failed to init state: %+v", err)
		return transDtoErrorToEcho(err)
	}
	result := map[string]string{"state": state}
	return ctx.JSON(http.StatusOK, result)
}

func (c *authController) Login(ctx echo.Context) error {
	// リクエストパース
	var req LoginRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		ctx.Logger().Infof("Failed to decode request: %+v", err)
		return echo.ErrBadRequest
	}
	defer ctx.Request().Body.Close()
	if err := validate.Struct(req); err != nil {
		ctx.Logger().Infof("Failed to validate request: %+v", err)
		return echo.ErrBadRequest
	}

	// リダイレクト先のドメインチェック
	redirectUrl, err := url.Parse(req.RedirectTo)
	if err != nil {
		ctx.Logger().Infof("Failed to parse specified redirect url: %+v", err)
		return echo.ErrBadRequest
	}
	if !config.IsAllowedDomain(redirectUrl.Host) {
		ctx.Logger().Infof("Not allowed domain: %s", redirectUrl.Host)
		return echo.ErrBadRequest
	}

	// ログイン処理
	if err := c.authUsecase.Login(ctx, &inputport.AuthLoginRequest{
		Username: req.Username,
		Password: req.Password,
		State:    req.State,
	}); err != nil {
		ctx.Logger().Infof("Failed to login: %+v", err)
		return transDtoErrorToEcho(err)
	}

	return ctx.Redirect(http.StatusSeeOther, req.RedirectTo)
}

func (c *authController) Logout(ctx echo.Context) error {
	return c.Logout(ctx)
}

func transDtoErrorToEcho(e error) error {
	switch e.(type) {
	case *dto.AuthBadRequestError:
		return echo.ErrBadRequest
	case *dto.AuthUnauthorizedError:
		return echo.ErrUnauthorized
	case *dto.AuthForbiddenError:
		return echo.ErrForbidden
	default:
		return echo.ErrInternalServerError
	}
}
