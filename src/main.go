package main

import (
	"log"
	"os"

	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/ai"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/bot"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/db"
)

// const DEBUG bool = true

func main() {

	func() {

		// Токен телеграм бота
		TELEGRAM_TOKEN, existsTT := os.LookupEnv("TELEGRAM_TOKEN")

		// // Токен OpenAI для ИИ
		AI_TOKEN, existsAT := os.LookupEnv("AI_TOKEN")

		// // Данные для подключения к базе данных
		DB_NAME, existsDN := os.LookupEnv("DB_NAME")
		DB_LOGIN, existsDL := os.LookupEnv("DB_LOGIN")
		DB_PASSWORD, existsDP := os.LookupEnv("DB_PASSWORD")

		// Название ИИ от OpenAI
		AI_MODEL := "text-davinci-003"

		// Тайм аут телеграм бота в секундах
		TELEGRAM_TIMEOUT := 60

		//////////////////////////////////////////////////
		if existsTT {
			log.Fatal("Отсутвует глобальная пременная с токеном телеги.")
		}
		if existsAT {
			log.Fatal("Отсутвует глобальная пременная с токеном OpenAI.")
		}
		if existsDN {
			log.Fatal("Отсутвует глобальная пременная с именем БД.")
		}
		if existsDL {
			log.Fatal("Отсутвует глобальная пременная с логином ДБ.")
		}
		if existsDP {
			log.Fatal("Отсутвует глобальная пременная с паролем от ДБ.")
		}

		// Инициализация компонентов
		bot.Init(TELEGRAM_TOKEN, TELEGRAM_TIMEOUT, false)
		ai.Init(AI_TOKEN, AI_MODEL)
		db.Init(DB_LOGIN, DB_PASSWORD, DB_NAME)

	}()

	defer db.Close()
	bot.Run()
}
