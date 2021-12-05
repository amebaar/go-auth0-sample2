package outputport

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type AuthClient interface {
	GetToken(ctx context.Context, username, password string) (*oauth2.Token, error)
	VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
}
