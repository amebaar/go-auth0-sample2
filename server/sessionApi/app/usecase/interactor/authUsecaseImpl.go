package interactor

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/server/sessionApi/app/usecase/dto"
	"go-auth0-sample2/server/sessionApi/app/usecase/inputport"
	"go-auth0-sample2/server/sessionApi/app/usecase/outputport"
	"math/rand"
)

type authUsecase struct {
	authClient outputport.AuthClient
}

func NewAuthUsecase(auth0client outputport.AuthClient) inputport.AuthUsecase {
	return &authUsecase{auth0client}
}

func (u *authUsecase) InitState(ctx echo.Context) (string, error) {
	state, err := generateRandomState()
	if err != nil {
		return "", &dto.AuthInternalServerError{
			BaseErr: err,
		}
	}

	// Save the state inside the session. [COOKIE]
	sess, _ := session.Get("session", ctx)
	sess.Values["state"] = state
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return "", &dto.AuthInternalServerError{
			BaseErr: err,
		}
	}
	return state, nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func (u *authUsecase) Login(ctx echo.Context, request *inputport.AuthLoginRequest) error {
	sess, err := session.Get("session", ctx)
	if err != nil || request.State != sess.Values["state"] {
		return &dto.AuthBadRequestError{
			BaseErr: fmt.Errorf(
				fmt.Sprintf("invalid state parameter %#v, %#v: %%w", request.State, sess.Values["state"]),
				err),
		}
	}

	token, err := u.authClient.GetToken(ctx.Request().Context(), request.Username, request.Password)
	if err != nil {
		return &dto.AuthUnauthorizedError{
			BaseErr: fmt.Errorf(
				"failed to convert an authorization code into a token: %w",
				err),
		}
	}

	idToken, err := u.authClient.VerifyIDToken(ctx.Request().Context(), token)
	if err != nil {
		return &dto.AuthInternalServerError{
			BaseErr: fmt.Errorf(
				"failed to verify id token: %w",
				err),
		}
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return &dto.AuthInternalServerError{
			BaseErr: fmt.Errorf(
				"failed get claims: %w",
				err),
		}
	}

	sess.Values["access_token"] = token.AccessToken
	sess.Values["profile"] = profile
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return &dto.AuthInternalServerError{
			BaseErr: fmt.Errorf(
				"failed to save session: %w",
				err),
		}
	}

	return nil
}
