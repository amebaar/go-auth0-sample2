package gateway

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
	"go-auth0-sample2/server/sessionApi/app/domain/repository"
	"math/rand"
)

type sessionRepository struct{}

func NewSessionRepository() repository.SessionRepository {
	return &sessionRepository{}
}

func (sr *sessionRepository) InitState(ctx echo.Context) (string, error) {
	generateRandomState := func() (string, error) {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		state := base64.StdEncoding.EncodeToString(b)
		return state, nil
	}

	state, err := generateRandomState()
	if err != nil {
		return "", err
	}

	// Save the state inside the session. [COOKIE]
	sess, _ := session.Get("session", ctx)
	sess.Values["state"] = state
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return "", err
	}
	return state, nil
}

func (sr *sessionRepository) IsValidState(ctx echo.Context, state string) bool {
	sess, err := session.Get("session", ctx)
	fmt.Printf("%v, %v", state, sess.Values["state"])
	return err == nil && state == sess.Values["state"]
}

func (sr *sessionRepository) Save(ctx echo.Context, userContext model.UserContext) error {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return err
	}

	var profile map[string]interface{}
	if err := userContext.LoadClaims(&profile); err != nil {
		return err
	}
	sess.Values["profile"] = profile
	sess.Values["access_token"] = userContext.AccessToken()

	return sess.Save(ctx.Request(), ctx.Response())
}

func (sr *sessionRepository) Discard(ctx echo.Context) error {
	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	sess.Values = map[interface{}]interface{}{}
	return sess.Save(ctx.Request(), ctx.Response())
}
