package model

type Role struct {
	id          string
	name        string
	permissions []*Permission
}

type Permission struct {
	operation  string
	resource   string
	identifier string
}
