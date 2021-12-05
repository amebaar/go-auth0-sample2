package main

import (
	"github.com/joho/godotenv"
	"go-auth0-sample2/server/sessionApi/app/adapter/web/controller"
	"go-auth0-sample2/server/sessionApi/app/infrastructure/auth0"
	"go-auth0-sample2/server/sessionApi/app/infrastructure/web"
	"go-auth0-sample2/server/sessionApi/app/usecase/interactor"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %+v", err)
	}

	authClient, err := auth0.New()
	if err != nil {
		panic(err)
	}
	authUseCase := interactor.NewAuthUsecase(authClient)
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
