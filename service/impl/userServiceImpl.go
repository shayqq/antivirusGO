package impl

import (
	"antivirus/model"
	"antivirus/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct{}

func (a *UserServiceImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (a *UserServiceImpl) Authenticate(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (a *UserServiceImpl) GetAll() ([]model.ApplicationUser, error) {
	return user.FindAll()
}
