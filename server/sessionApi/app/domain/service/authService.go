package service

import (
	"context"
	"go-auth0-sample2/server/sessionApi/app/domain/model"
	"go-auth0-sample2/server/sessionApi/app/domain/repository"
)

type AuthService interface {
	FetchUserContextByUserCredential(ctx context.Context, cred model.UserCredential) (model.UserContext, error)
}

type authService struct {
	tokenRepository repository.TokenRepository
}

func NewAuthService(tokenRepository repository.TokenRepository) AuthService {
	return &authService{tokenRepository}
}

func (as *authService) FetchUserContextByUserCredential(ctx context.Context, cred model.UserCredential) (model.UserContext, error) {
	token, err := as.tokenRepository.GetToken(ctx, cred)
	if err != nil {
		return nil, err
	}
	idToken, err := as.tokenRepository.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return model.NewUserContext(token, idToken), nil
}
