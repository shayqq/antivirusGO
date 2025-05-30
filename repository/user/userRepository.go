package user

import (
	"antivirus/database"
	"antivirus/model"
	"database/sql"
)

func FindByEmail(postEmail string) (*model.ApplicationUser, error) {
	applicationUser := model.ApplicationUser{}
	query := `SELECT username, email, password, role FROM users WHERE email=$1`
	err := database.DB.QueryRow(query, postEmail).Scan(&applicationUser.Username, &applicationUser.Email,
		&applicationUser.Password, &applicationUser.Role)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return &applicationUser, nil
}

func FindById(id int) (*model.ApplicationUser, error) {
	applicationUser := model.ApplicationUser{}
	query := `SELECT username, email, password, role FROM users WHERE id=$1`
	err := database.DB.QueryRow(query, id).Scan(&applicationUser.Username, &applicationUser.Email,
		&applicationUser.Password, &applicationUser.Role)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return &applicationUser, nil
}

func FindAll() ([]model.ApplicationUser, error) {
	query := `SELECT username, email, role FROM users`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.ApplicationUser
	for rows.Next() {
		var user model.ApplicationUser
		err := rows.Scan(&user.Username, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}
