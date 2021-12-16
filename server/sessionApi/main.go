package main

import (
	"github.com/joho/godotenv"
	"go-auth0-sample2/sdk/authman"
	"go-auth0-sample2/server/sessionApi/app/adapter/gateway"
	"go-auth0-sample2/server/sessionApi/app/adapter/web/controller"
	"go-auth0-sample2/server/sessionApi/app/domain/service"
	"go-auth0-sample2/server/sessionApi/app/infrastructure/auth"
	"go-auth0-sample2/server/sessionApi/app/infrastructure/web"
	"go-auth0-sample2/server/sessionApi/app/usecase/interactor"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %+v", err)
	}

	// TODO("DIフレームワーク導入。uber/diとか")
	authMan, err := authman.GetAuthMan()
	if err != nil {
		panic(err)
	}
	authClient := auth.GetClient()
	tokenRepo, err := gateway.NewTokenRepository(authClient)
	if err != nil {
		panic(err)
	}
	authSvc := service.NewAuthService(tokenRepo)
	sessRepo := gateway.NewSessionRepository()
	authUseCase := interactor.NewAuthUsecase(authMan, authSvc, sessRepo, tokenRepo)
	authController := controller.NewAuthController(authUseCase)

	web.Start(authController)
	/*
		cl, err := auth0.New()
		if err != nil {
			panic(err)
		}
		tk, err := cl.GetToken(context.Background(), "contact.maebara@gmail.com", "Sample1234!")
		if err != nil {
			panic(err)
		}

		fmt.Println(tk)
	*/
}
