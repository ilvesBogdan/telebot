package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

/*
Init - функция инициализации базы данных.

Параметры:

	login - логин пользователя;
	password - пароль;
	dataBaseName - имя базы данных.

Функция создает соединение с указанной базой данных, средствами sqlx.
*/
func Init(login string, password string, dataBaseName string) {
	var err error

	DB, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("sslmode=disable user=%v password=%v dbname=%v", login, password, dataBaseName))
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}

	log.Print("Подключено к базе данных. Номер подключения: ", DB.Stats().Idle)
}

/*
Close - функция закрытия соединения с базой данных.

Функция закрывает соединение с базой данных, средствами sqlx.
*/
func Close() {
	err := DB.Close()
	if err != nil {
		log.Fatal("Не удалось закрыть соединение с базой данных: ", err)
	}
}
