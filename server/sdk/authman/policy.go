package authman

import (
	"go-auth0-sample2/sdk/authman/infra"
	"os"
	"strings"
)

type policyManager struct {
	client     infra.Auth0Client
	identifier string
	roles      []*role
}

type permission struct {
	operation  string
	resource   string
	identifier string
}

type role struct {
	id          string
	name        string
	permissions []*permission
}

func (m *policyManager) refresh() error {
	roles, err := m.client.GetRoles()
	if err != nil {
		return err
	}

	result := make([]*role, 0, len(roles))

	for _, r := range roles {
		permissions := make([]*permission, 0, len(r.Permissions))
		for _, p := range r.Permissions {
			ops := strings.Split(p.Name, ":")
			if len(ops) != 2 {
				continue
			}
			permissions = append(permissions, &permission{
				ops[0], ops[1], p.Identifier,
			})
		}
		result = append(result, &role{
			id:          r.Id,
			name:        r.Name,
			permissions: permissions,
		})
	}
	m.roles = result
	return nil
}

func (m *policyManager) hasPermission(role string, operation string, resource string) bool {
	if m.roles == nil {
		if m.refresh() != nil {
			return false
		}
	}

	for _, r := range m.roles {
		if r.name == role {
			for _, p := range r.permissions {
				if p.operation == operation && p.resource == resource && p.identifier == m.identifier {
					return true
				}
			}
		}
	}
	return false
}

func newPolicyManager(cl infra.Auth0Client) *policyManager {
	identifier := os.Getenv("AUTH0_API_IDENTIFIER")
	pm := &policyManager{
		cl, identifier, nil,
	}
	_ = pm.refresh()
	return pm
}
