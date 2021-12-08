package model

type UserCredential interface {
	GetName() string
	GetPassword() string
}

type userCredential struct {
	name     *userName
	password *userPassword
}

type userName struct {
	name string
}

func UserName(name string) *userName {
	// TODO("add verification")
	return &userName{name}
}

type userPassword struct {
	password string
}

func UserPassword(password string) *userPassword {
	// TODO("add verification")
	return &userPassword{password}
}

func NewUserCredential(name *userName, password *userPassword) UserCredential {
	return &userCredential{name, password}
}

func (uc *userCredential) GetName() string {
	return uc.name.name
}

func (uc *userCredential) GetPassword() string {
	return uc.password.password
}
