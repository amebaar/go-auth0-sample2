package auth

import (
	"gopkg.in/auth0.v5/management"
	"os"
	"sync"
)

type Client interface {
	CreateUser(req *CreateUserRequest) error
}

type client struct {
	client *management.Management
}

var once sync.Once
var instance Client

func GetClient() Client {
	once.Do(func() {
		domain := os.Getenv("AUTH0_DOMAIN")
		id := os.Getenv("AUTH0_CLIENT_ID")
		secret := os.Getenv("AUTH0_CLIENT_SECRET")
		m, err := management.New(domain, management.WithClientCredentials(id, secret))
		if err != nil {
			panic(err)
		}
		instance = &client{
			m,
		}
	})
	return instance
}

type CreateUserRequest struct {
	Email      string
	Password   string
	Connection string
	UserName   *string
}

func (c *client) CreateUser(req *CreateUserRequest) error {
	toPtrBool := func(b bool) *bool {
		return &b
	}
	user := &management.User{
		Connection:    &req.Connection,
		Email:         &req.Email,
		Username:      req.UserName,
		Password:      &req.Password,
		EmailVerified: toPtrBool(false),
		VerifyEmail:   toPtrBool(true),
	}
	return c.client.User.Create(user)
}
