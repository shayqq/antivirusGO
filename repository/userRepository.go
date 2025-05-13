package repository

import (
	"antivirus/database"
	"antivirus/model"
	"database/sql"
)

func FindByEmail(postEmail string) *model.ApplicationUser {
	applicationUser := model.ApplicationUser{}
	query := `SELECT username, email, password, role FROM users WHERE email=$1`
	err := database.DB.QueryRow(query, postEmail).Scan(&applicationUser.Username, &applicationUser.Email,
		&applicationUser.Password, &applicationUser.Role)
	if err == sql.ErrNoRows {
		return nil
	}
	return &applicationUser
}
