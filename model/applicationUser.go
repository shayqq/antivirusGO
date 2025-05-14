package model

type ApplicationUser struct {
	Username string
	Email    string
	Password string `json:"-"`
	Role     int8
}
