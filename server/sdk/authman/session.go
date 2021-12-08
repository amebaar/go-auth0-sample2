package authman

import (
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/sdk/authman/model"
)

const (
	RoleNamespace = "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"
	CompanyIdKey  = "cid"
	NameKey       = "name"
	StateKey      = "state"
)

type sessionManager struct{}

func (sr *sessionManager) loadUserProfile(ctx echo.Context) (model.UserProfile, error) {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return nil, err
	}

	profile := sess.Values["profile"].(map[string]interface{})
	if profile == nil {
		return nil, fmt.Errorf("failed to fetch user profile from session")
	}

	roles, ok := profile[RoleNamespace].([]string)
	if !ok {
		return nil, fmt.Errorf("invalid profile: roles is %+v", roles)
	}
	var cidPtr *int
	cid, ok := profile[CompanyIdKey].(int)
	if ok {
		cidPtr = &cid
	}
	name, ok := profile[NameKey].(string)
	if !ok {
		return nil, fmt.Errorf("invalid profile: name is %+v", name)
	}
	state, ok := profile[StateKey].(string)
	if !ok {
		return nil, fmt.Errorf("invalid profile: state is %+v", state)
	}

	return model.NewUserProfile(
		cidPtr,
		model.UserName(name),
		model.UserRoles(roles),
		state,
	), nil
}

func newSessionManager() *sessionManager {
	return &sessionManager{}
}
