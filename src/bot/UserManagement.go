package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/db"
)

func getUserByID(userID int64) (user db.User) {
	user = db.GetUserByTelegramID(userID)
	user.TrimSpace()
	return
}

func registerNewUser(name string, user *tgbotapi.User) {
	db.SetNewUser(name, user.UserName, user.FirstName, user.LastName, user.ID)
}
