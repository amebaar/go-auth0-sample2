package infra

import (
	"gopkg.in/auth0.v5"
	"gopkg.in/auth0.v5/management"
	"os"
	"sync"
)

type Auth0Client interface {
	GetRoles() ([]*role, error)
}

type auth0Client struct {
	client *management.Management
}

var once sync.Once
var instance Auth0Client

func GetAuth0Client() Auth0Client {
	once.Do(func() {
		domain := os.Getenv("AUTH0_DOMAIN")
		id := os.Getenv("AUTH0_CLIENT_ID")
		secret := os.Getenv("AUTH0_CLIENT_SECRET")
		m, err := management.New(domain, management.WithClientCredentials(id, secret))
		if err != nil {
			panic(err)
		}
		instance = &auth0Client{
			m,
		}
	})
	return instance
}

type permission struct {
	Name       string
	Identifier string
}

type role struct {
	Id          string
	Name        string
	Permissions []*permission
}

func (c *auth0Client) GetRoles() ([]*role, error) {
	list, err := c.client.Role.List()
	if err != nil {
		return nil, err
	}

	result := make([]*role, 0, len(list.Roles))

	for _, v := range list.Roles {
		p, err := c.getPermissions(auth0.StringValue(v.ID))
		if err != nil {
			return nil, err
		}
		result = append(result, &role{
			auth0.StringValue(v.ID), auth0.StringValue(v.Name), p,
		})
	}
	return result, nil
}

func (c *auth0Client) getPermissions(roleId string) ([]*permission, error) {
	list, err := c.client.Role.Permissions(roleId)
	if err != nil {
		return nil, err
	}

	result := make([]*permission, 0, len(list.Permissions))
	for _, v := range list.Permissions {
		result = append(result, &permission{
			auth0.StringValue(v.Name), auth0.StringValue(v.ResourceServerIdentifier),
		})
	}
	return result, nil
}
