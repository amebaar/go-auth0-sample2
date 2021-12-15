package repository

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
	"golang.org/x/oauth2"
	"net/url"
)

type TokenRepository interface {
	GetToken(ctx context.Context, cred model.UserCredential) (*oauth2.Token, error)
	GetAuthCodeUrl(ctx context.Context, state string, conn string, redirectTo *url.URL) (string, error)
	GetTokenByCode(ctx context.Context, code string) (*oauth2.Token, error)
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}
