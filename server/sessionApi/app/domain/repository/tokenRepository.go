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
	CreateUser(ctx context.Context, cred model.UserCredential) error                                    // TODO("move to another package")
	GetAuthCodeUrl(ctx context.Context, state string, conn string, redirectTo *url.URL) (string, error) // TODO("move to another package")
	GetTokenByCode(ctx context.Context, code string) (*oauth2.Token, error)
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}
