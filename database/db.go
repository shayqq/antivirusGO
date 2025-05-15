package database

import (
	"antivirus/model"
	"antivirus/textcolor"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := "host=localhost user=" + os.Getenv("DBUSERNAME") + " password=" + os.Getenv("DBPASSWORD") + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(textcolor.RedErrorText("Не удалось открыть базу данных "), err)
	}
	query := "CREATE DATABASE " + os.Getenv("DBNAME")
	_, err = db.Exec(query)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatal(textcolor.RedErrorText("Что-то пошло не так "), err)
			db.Close()
		} else {
			db.Close()
			goto Next
		}
	} else {
		fmt.Println(textcolor.GreenSuccessText("База данных успешно создана!"))
		db.Close()
	}
Next:
	dsnMain := "host=localhost user=" + os.Getenv("DBUSERNAME") + " password=" + os.Getenv("DBPASSWORD") + " dbname=" +
		os.Getenv("DBNAME") + " sslmode=disable"
	mainDb, err := sql.Open("postgres", dsnMain)
	if err != nil {
		log.Fatal(textcolor.RedErrorText("Не удалось открыть базу данных "), err)
	}
	err = mainDb.Ping()
	if err != nil {
		log.Fatal(textcolor.RedErrorText("Не удалось подключиться к базе данных "), err)
	}
	DB = mainDb
	fmt.Println(textcolor.GreenSuccessText("Подключение к базе данных успешно!"))
	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role SMALLINT NOT NULL DEFAULT 0 CHECK (role IN (0, 1))
	)
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(textcolor.RedErrorText("Не удалось создать таблицы "), err)
	}
	fmt.Println(textcolor.GreenSuccessText("Таблицы успешно созданы!"))
}

func Save(object any) string {
	switch v := object.(type) {
	case model.ApplicationUser:
		query := `INSERT INTO users (username, email, password, role) values ($1, $2, $3, $4)`
		_, err := DB.Exec(query, v.Username, v.Email, v.Password, v.Role)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "username") {
				return "Пользователь с таким логином уже существует"
			} else if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "email") {
				return "Пользователь с таким email уже существует"
			}
			return "Ошибка сервера"
		}
		return ""
	default:
		return "Ошибка сервера"
	}
}

func Update(object any, flag any) string {
	switch v := object.(type) {
	case model.ApplicationUser:
		query := `UPDATE users SET username=$1, email=$2, password=$3, role=$4 WHERE id=$5`
		_, err := DB.Exec(query, v.Username, v.Email, v.Password, v.Role, flag)
		if err != nil {
			return "Ошибка сервера"
		}
		return ""
	default:
		return "Ошибка сервера"
	}
}

func Delete(object any, flag any) string {
	switch object.(type) {
	case model.ApplicationUser:
		query := `DELETE FROM users WHERE id=$1`
		_, err := DB.Exec(query, flag)
		if err != nil {
			return "Ошибка сервера"
		}
		return ""
	default:
		return "Ошибка сервера"
	}
}
