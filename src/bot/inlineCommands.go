package bot

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/db"
)

func SendCallbackMessage(callbackQuery *tgbotapi.CallbackQuery, text string) {
	callback := tgbotapi.NewCallback(callbackQuery.ID, text)
	if _, err := BOT.Request(callback); err != nil {
		log.Println("Ошибка отправки callback'а: ", err)
	}
}

func callbackInlineButton(update tgbotapi.Update) {

	switch update.CallbackQuery.Data {
	case "close":
		// Если нажата кнопка закрытия, удаляем сообщение
		removeMessageByCallbackQuery(update.CallbackQuery)
	case "stop-generate":
		// todo Реализовать остановку генерации текста
	default:
		// Проверяем на наличие действия для кнопок
		// изменения контекста
		if callbackMessageContextButtons(update.CallbackQuery) {
			return
		}
		// todo исправить текст
		SendCallbackMessage(update.CallbackQuery, "Ошибка сервера")
	}
}

func callbackMessageContextButtons(callbackQuery *tgbotapi.CallbackQuery) bool {
	resultSplit := strings.Split(callbackQuery.Data, ":")
	if len(resultSplit) != 2 {
		return false
	}
	switch resultSplit[0] {
	case "ctxCnfgClear":
		// Нажата кнопка "отчистить контекст"
		removeMessageByCallbackQuery(callbackQuery)
		contextID := getContextIDWithCallback(resultSplit)
		if contextID > 0 {
			db.SetNewContextText(contextID, "")
			// todo Добавить закрытие списка конекстов, если он открыт
			// todo Исправить текст
			SendCallbackMessage(callbackQuery, "Контекст отчищен")
			return true
		}
	case "ctxCnfgRename":
		// Нажата кнопка переименовать контекст
		removeMessageByCallbackQuery(callbackQuery)
		contextID := getContextIDWithCallback(resultSplit)
		if contextID > 0 && setExpectedMessageTypeForUser(3, &callbackQuery.From.ID) {
			// todo Испраавить текст сообщения
			quickSendTextMessage(callbackQuery.From.ID, "Введите новое имя для контекста")
			return true
		}
		return true
	case "ctxCnfgRemove":
		removeMessageByCallbackQuery(callbackQuery)
		contextID := getContextIDWithCallback(resultSplit)
		if contextID > 0 {

			if db.RemoveContextByID(contextID) {
				SendCallbackMessage(callbackQuery, "Контекст успешно удален")
			} else {
				SendCallbackMessage(callbackQuery, "Ошибка удаления контекста")
			}

			// Если удаленый контекст был текущи, то переключаем его
			if getUserByID(callbackQuery.From.ID).Context == contextID {
				withoutContext(&callbackQuery.From.ID)
			}
			return true
		}
	}
	return false
}

func removeMessageByCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	BOT.Send(tgbotapi.NewDeleteMessage(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
	))
}

func getContextIDWithCallback(resultSplit []string) int {
	contextID, err := strconv.Atoi(resultSplit[1])
	if err != nil {
		log.Printf("Ошибка при парсинге числа из '%v': %v", resultSplit[1], err)
		return 0
	}
	return contextID
}
