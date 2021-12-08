package authman

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go-auth0-sample2/sdk/authman/infra"
	"go-auth0-sample2/sdk/authman/model"
	"sync"
)

type AuthMan interface {
	ReloadPolicy() error
	HasPermission(ctx echo.Context, operation string, resource string) bool
	LoadProfile(ctx echo.Context) (model.UserProfile, error)
}

type authMan struct {
	pm *policyManager
	sm *sessionManager
}

var once sync.Once
var instance AuthMan

func GetAuthMan() (AuthMan, error) {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
		instance = &authMan{newPolicyManager(infra.GetAuth0Client()), newSessionManager()}
	})
	return instance, nil
}

func (am *authMan) ReloadPolicy() error {
	return am.pm.refresh()
}

func (am *authMan) HasPermission(ctx echo.Context, operation string, resource string) bool {
	p, err := am.sm.loadUserProfile(ctx)
	if err != nil {
		return false
	}

	for _, r := range p.GetRoles() {
		if am.pm.hasPermission(r, operation, resource) {
			return true
		}
	}

	return false
}

func (am *authMan) LoadProfile(ctx echo.Context) (model.UserProfile, error) {
	return am.sm.loadUserProfile(ctx)
}
