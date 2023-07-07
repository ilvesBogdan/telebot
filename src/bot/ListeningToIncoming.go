package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handlerIncomingMessages(update tgbotapi.Update) {

	user := getUserByID(update.Message.From.ID)
	text := setLanguage(update.Message.From.LanguageCode)

	if user.ID == 0 {

		// Пользователь не зарегистрирован

		if name, check := checkNewUserInWaiting(&update.Message.From.UserName); check {
			registerNewUser(name, update.Message.From)
			quickSendTextMessage(update.Message.Chat.ID, "Привет!")
			return
		}

		quickSendTextMessage(update.Message.Chat.ID, "Бан!")

	} else {

		if commands(&text, &update) {
			return
		}

		switch user.MessageType {
		case 0:
			// Пользователь не должен отправлять сообщения
			// todo Исправить текст
			quickSendTextMessage(update.Message.Chat.ID, "Бан")
		case 1:
			// Пользователь отправил сообщение к ИИ
			if user.Context == 0 {
				go sendingMessagesFromAI("", update.Message.Text, *update.Message)
			} else {
				context := getContext(&user.Context)
				go sendingMessagesFromAI(context.Text, update.Message.Text, *update.Message)
				go setNewContext(context.Text, update.Message.Text, context.Id)
				go writeMessageInDB(context.Id, update.Message.Text)
			}
		case 2:
			// Пользователь вводит название нового контекста.
			if checkNameContext(update.Message) {
				return
			}
			createNewContextAndSelect(update.Message.Text, update.Message)
			if !setExpectedMessageTypeForUser(1, &update.Message.From.ID) {
				// todo Исправить текст сообщения
				quickSendTextMessage(update.Message.From.ID, "Ошибка сервера")
			}
		case 3:
			// Пользователь переименовывает существующий контекст
			if checkNameContext(update.Message) {
				return
			}
			renameContext(update.Message, &user.Context)
			if !setExpectedMessageTypeForUser(1, &update.Message.From.ID) {
				// todo Исправить текст сообщения
				quickSendTextMessage(update.Message.From.ID, "Ошибка сервера")
			}
		}
	}
}
