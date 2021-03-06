package model

import "fmt"

type UserProfile interface {
	GetCompanyId() (int, error)
	GetName() string
	GetRoles() []string
}

type userProfile struct {
	companyId int
	name      *userName
	roles     *userRoles
	// add some user metadata which must be stored in session
}

type consumerProfile struct {
	name  *userName
	roles *userRoles
	// add some user metadata which must be stored in session
}

type userName struct {
	name string
}

func UserName(name string) *userName {
	// TODO("add verification")
	return &userName{name}
}

type userRoles struct {
	roles []string
}

func UserRoles(roles []string) *userRoles {
	// TODO("add verification")
	return &userRoles{roles}
}

func NewUserProfile(
	companyId *int,
	name *userName,
	roles *userRoles,
) UserProfile {
	if companyId == nil {
		return &consumerProfile{
			name, roles,
		}
	} else {
		return &userProfile{
			*companyId, name, roles,
		}
	}
}

func (up *userProfile) GetCompanyId() (int, error) {
	return up.companyId, nil
}
func (up *userProfile) GetName() string {
	return up.name.name
}
func (up *userProfile) GetRoles() []string {
	return up.roles.roles
}

func (cp *consumerProfile) GetCompanyId() (int, error) {
	return 0, fmt.Errorf("user type is Consumer")
}
func (cp *consumerProfile) GetName() string {
	return cp.name.name
}
func (cp *consumerProfile) GetRoles() []string {
	return cp.roles.roles
}
