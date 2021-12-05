package auth0

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"go-auth0-sample2/server/sessionApi/app/usecase/outputport"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type client struct {
	*oidc.Provider
	oauth2.Config
	*http.Client
}

func New() (outputport.AuthClient, error) {
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
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	httpClient := &http.Client{} // TODO("コネクションプールの設定など")

	return &client{
		provider,
		conf,
		httpClient,
	}, nil
}

func (c *client) GetToken(ctx context.Context, username, password string) (*oauth2.Token, error) {
	newCtx := context.WithValue(ctx, oauth2.HTTPClient, c.Client) // 使用するHTTP Clientを指定
	return c.PasswordCredentialsToken(newCtx, username, password)
}

func (c *client) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: c.ClientID,
	}

	return c.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
