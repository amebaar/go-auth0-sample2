package web

import (
	"encoding/gob"
	"go-auth0-sample2/server/sessionApi/app/adapter/web/controller"
)

func Start(
	authController controller.AuthController,
) {
	gob.Register(map[string]interface{}{})
	router := newRouter(authController)
	router.Logger.Fatal(router.Start(":8080"))
}
