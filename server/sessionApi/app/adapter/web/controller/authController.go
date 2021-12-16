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
	SignUp(ctx echo.Context) error
	Login(ctx echo.Context) error
	SocialLogin(ctx echo.Context) error
	Callback(ctx echo.Context) error
	Logout(ctx echo.Context) error
}

type authController struct {
	authUsecase inputport.AuthUsecase
}

func NewAuthController(authUsecase inputport.AuthUsecase) AuthController {
	validate = validator.New()
	return &authController{authUsecase}
}

type SignUpRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RedirectTo string `json:"redirect_to" validate:"required,url"`
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

func (c *authController) SignUp(ctx echo.Context) error {
	// リクエストパース
	var req SignUpRequest
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

	// サインアップ処理
	if err := c.authUsecase.SignUp(ctx, &inputport.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		ctx.Logger().Infof("Failed to sign up: %+v", err)
		return transDtoErrorToEcho(err)
	}

	return ctx.Redirect(http.StatusSeeOther, req.RedirectTo)
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

func (c *authController) SocialLogin(ctx echo.Context) error {
	// リクエストパース
	connection := ctx.QueryParam("connection")
	state := ctx.QueryParam("state")
	redirectTo := ctx.QueryParam("redirect")
	if connection == "" || state == "" {
		ctx.Logger().Infof("Invalid query: %+v", ctx.QueryParams())
		return echo.ErrBadRequest
	}

	// リダイレクト先のドメインチェック
	redirectUrl, err := url.Parse(redirectTo)
	if err != nil {
		ctx.Logger().Infof("Failed to parse specified redirect url: %+v", err)
		return echo.ErrBadRequest
	}
	if !config.IsAllowedDomain(redirectUrl.Host) {
		ctx.Logger().Infof("Not allowed domain: %s", redirectUrl.Host)
		return echo.ErrBadRequest
	}

	// ログイン処理
	redirect, err := c.authUsecase.SocialLogin(ctx, &inputport.SocialLoginRequest{
		Connection: connection,
		State:      state,
		RedirectTo: redirectTo,
	})
	if err != nil {
		ctx.Logger().Infof("Failed to login: %+v", err)
		return transDtoErrorToEcho(err)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, redirect)
}

func (c *authController) Callback(ctx echo.Context) error {
	// リクエストパース
	src := ctx.QueryParam("src")
	code := ctx.QueryParam("code")
	state := ctx.QueryParam("state")
	if code == "" || state == "" {
		ctx.Logger().Infof("Invalid query: %+v", ctx.QueryParams())
		return echo.ErrBadRequest
	}

	// リダイレクト先のドメインチェック
	redirectUrl, err := url.Parse(src)
	if err != nil {
		ctx.Logger().Infof("Failed to parse specified redirect url: %+v", err)
		return echo.ErrBadRequest
	}
	if !config.IsAllowedDomain(redirectUrl.Host) {
		ctx.Logger().Infof("Not allowed domain: %s", redirectUrl.Host)
		return echo.ErrBadRequest
	}

	// セッション更新処理
	if err := c.authUsecase.UpdateSession(ctx, &inputport.UpdateSessionRequest{
		State: state,
		Code:  code,
	}); err != nil {
		ctx.Logger().Infof("Failed to login: %+v", err)
		return transDtoErrorToEcho(err)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, src)
}

func (c *authController) Logout(ctx echo.Context) error {
	if err := c.authUsecase.Logout(ctx); err != nil {
		return echo.ErrInternalServerError
	}
	return ctx.NoContent(http.StatusOK)
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
