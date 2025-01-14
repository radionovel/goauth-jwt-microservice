package model

type NewUserDTO struct {
	UserProvider int
	Username     string
	Password     string
}

type User struct {
	ID       int
	Username string
}
