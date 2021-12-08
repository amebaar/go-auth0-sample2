package model

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type UserContext interface {
	AccessToken() string
	LoadClaims(profile interface{}) error
}

type userContext struct {
	*oauth2.Token
	*oidc.IDToken
}

func NewUserContext(token *oauth2.Token, idToken *oidc.IDToken) UserContext {
	return &userContext{token, idToken}
}

func (uc *userContext) AccessToken() string {
	return uc.Token.AccessToken
}

func (uc *userContext) LoadClaims(profile interface{}) error {
	return uc.IDToken.Claims(profile)
}
