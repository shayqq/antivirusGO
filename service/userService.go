package service

type UserService interface {
	HashPassword(password string) (string, error)
	Authenticate(hashedPassword string, plainPassword string) bool
}
