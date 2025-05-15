package model

type ApplicationUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int8   `json:"role"`
}
