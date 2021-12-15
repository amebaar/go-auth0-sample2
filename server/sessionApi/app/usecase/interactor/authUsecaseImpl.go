package interactor

import (
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/sdk/authman"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
	"go-auth0-sample2/server/sessionApi/app/domain/repository"
	"go-auth0-sample2/server/sessionApi/app/domain/service"
	"go-auth0-sample2/server/sessionApi/app/usecase/dto"
	"go-auth0-sample2/server/sessionApi/app/usecase/inputport"
	"net/url"
)

type authUsecase struct {
	authMan         authman.AuthMan
	authService     service.AuthService
	sessRepository  repository.SessionRepository
	tokenRepository repository.TokenRepository
}

func NewAuthUsecase(
	authMan authman.AuthMan,
	authService service.AuthService,
	sessRepository repository.SessionRepository,
	tokenRepository repository.TokenRepository,
) inputport.AuthUsecase {
	return &authUsecase{authMan, authService, sessRepository, tokenRepository}
}

func (u *authUsecase) InitState(ctx echo.Context) (string, error) {
	return u.sessRepository.InitState(ctx)
}

func (u *authUsecase) Session(ctx echo.Context) (string, error) {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return "", &dto.AuthUnauthorizedError{
			BaseErr: fmt.Errorf(
				"failed to get session: %w",
				err),
		}
	}
	profile := sess.Values["profile"]
	if profile == nil {
		return "", &dto.AuthUnauthorizedError{
			BaseErr: fmt.Errorf(
				"failed to get profile: %w",
				err),
		}
	}

	return fmt.Sprint(profile), nil
}

func (u *authUsecase) Login(ctx echo.Context, request *inputport.AuthLoginRequest) error {
	if !u.sessRepository.IsValidState(ctx, request.State) {
		return fmt.Errorf("invalid state")
	}

	name := model.UserName(request.Username)
	pass := model.UserPassword(request.Password)

	userContext, err := u.authService.FetchUserContextByUserCredential(
		ctx.Request().Context(),
		model.NewUserCredential(name, pass),
	)
	if err != nil {
		return err
	}

	return u.sessRepository.Save(ctx, userContext)
}

func (u *authUsecase) SocialLogin(ctx echo.Context, request *inputport.SocialLoginRequest) (string, error) {
	if !u.sessRepository.IsValidState(ctx, request.State) {
		return "", fmt.Errorf("invalid state")
	}

	redirectTo, err := url.Parse(request.RedirectTo)
	if err != nil {
		return "", err
	}

	return u.tokenRepository.GetAuthCodeUrl(ctx.Request().Context(), request.State, request.Connection, redirectTo)
}

func (u *authUsecase) UpdateSession(ctx echo.Context, request *inputport.UpdateSessionRequest) error {
	if !u.sessRepository.IsValidState(ctx, request.State) {
		return fmt.Errorf("invalid state")
	}

	userContext, err := u.authService.FetchUserContextByCode(
		ctx.Request().Context(),
		request.Code,
	)
	if err != nil {
		return err
	}

	return u.sessRepository.Save(ctx, userContext)
}

func (u *authUsecase) Logout(ctx echo.Context) error {
	return u.sessRepository.Discard(ctx)
}
