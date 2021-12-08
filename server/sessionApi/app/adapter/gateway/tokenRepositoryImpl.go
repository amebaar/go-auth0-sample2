package gateway

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
	"go-auth0-sample2/server/sessionApi/app/domain/repository"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type tokenRepository struct {
	*oidc.Provider
	oauth2.Config
	*http.Client
}

func NewTokenRepository() (repository.TokenRepository, error) {
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

	return &tokenRepository{
		provider,
		conf,
		httpClient,
	}, nil
}

func (ar *tokenRepository) GetToken(ctx context.Context, cred model.UserCredential) (*oauth2.Token, error) {
	newCtx := context.WithValue(ctx, oauth2.HTTPClient, ar.Client) // 使用するHTTP Clientを指定
	return ar.PasswordCredentialsToken(newCtx, cred.GetName(), cred.GetPassword())
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
