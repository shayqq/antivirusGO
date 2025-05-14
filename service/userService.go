package service

import "antivirus/model"

type UserService interface {
	HashPassword(password string) (string, error)
	Authenticate(hashedPassword string, plainPassword string) bool
	GetAll() ([]model.ApplicationUser, error)
}
