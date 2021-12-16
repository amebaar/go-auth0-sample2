package gateway

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
	"go-auth0-sample2/server/sessionApi/app/domain/repository"
	"go-auth0-sample2/server/sessionApi/app/infrastructure/auth"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
)

type tokenRepository struct {
	*oidc.Provider
	oauth2.Config
	*http.Client
	authClient auth.Client
}

func NewTokenRepository(authClient auth.Client) (repository.TokenRepository, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(
		ctx,
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	httpClient := &http.Client{} // TODO("コネクションプールの設定など")

	return &tokenRepository{
		provider,
		conf,
		httpClient,
		authClient,
	}, nil
}

func (ar *tokenRepository) GetToken(ctx context.Context, cred model.UserCredential) (*oauth2.Token, error) {
	newCtx := context.WithValue(ctx, oauth2.HTTPClient, ar.Client) // 使用するHTTP Clientを指定
	return ar.PasswordCredentialsToken(newCtx, cred.GetName(), cred.GetPassword())
}

func (ar *tokenRepository) CreateUser(ctx context.Context, cred model.UserCredential) error {
	return ar.authClient.CreateUser(&auth.CreateUserRequest{
		Email:      cred.GetName(),
		Password:   cred.GetPassword(),
		Connection: "Username-Password-Authentication", // TODO
		UserName:   nil,
	})
}

func (ar *tokenRepository) GetAuthCodeUrl(ctx context.Context, state string, conn string, redirectTo *url.URL) (string, error) {
	redirectUrl := oauth2.SetAuthURLParam("redirect_uri", ar.Config.RedirectURL+"?src="+redirectTo.String())
	connection := oauth2.SetAuthURLParam("connection", conn)
	return ar.AuthCodeURL(state, redirectUrl, connection), nil
}

func (ar *tokenRepository) GetTokenByCode(ctx context.Context, code string) (*oauth2.Token, error) {
	newCtx := context.WithValue(ctx, oauth2.HTTPClient, ar.Client) // 使用するHTTP Clientを指定
	return ar.Exchange(newCtx, code)
}

func (ar *tokenRepository) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: ar.ClientID,
	}

	return ar.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
