package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/db"
)

func quickSendTextMessage(chatID int64, text string) (response tgbotapi.Message) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	response, err := BOT.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return
}

func userIsAdmin(user *tgbotapi.User) bool {
	// todo закончить
	return true
}

func setExpectedMessageTypeForUser(typeMsg int, userID *int64) bool {
	return db.SetExpectedMessageType(typeMsg, userID)
}

func sendConfigContextMessage(message *tgbotapi.Message, context *db.Context) {
	text := setLanguage(message.From.LanguageCode)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(
		text.msg("ContextCnfgMsg"), context.Name, context.GetDateStr()),
	)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text.msg("Clear"), fmt.Sprintf("ctxCnfgClear:%v", context.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text.msg("Rename"), fmt.Sprintf("ctxCnfgRename:%v", context.Id)),
			tgbotapi.NewInlineKeyboardButtonData(text.msg("Remove"), fmt.Sprintf("ctxCnfgRemove:%v", context.Id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			// todo Изменить текст сообщения
			tgbotapi.NewInlineKeyboardButtonData("Закрыть", "close"),
		),
	)

	if _, err := BOT.Send(msg); err != nil {
		log.Println(err)
	}
}

func showMenuContexts(message *tgbotapi.Message) {

	contexts, err := db.GetContextsByTid(message.From.ID)
	if nil != err {
		log.Panic("Неудалось получить данные контекстов из ДБ:", err)
	}

	if len(contexts) > 0 {
		var buttons [][]tgbotapi.KeyboardButton
		var btnText string

		for _, ct := range contexts {
			ct.TrimSpace()

			// Проверка на то, является ли этот контекст текущим
			if ct.Select {
				btnText = fmt.Sprintf("✅ %v", ct.Name)
			} else {
				btnText = fmt.Sprintf("💬 %v", ct.Name)
			}

			buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(btnText),
			))
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "📋")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)
		if _, err := BOT.Send(msg); err != nil {
			log.Println(err)
		}
	} else {
		// todo Исправить текст сообщения
		quickSendTextMessage(message.Chat.ID, "Нет контекстов.")
	}
	BOT.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
}

func hideMenuContexts(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "😎")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	res, err := BOT.Send(msg)
	if err != nil {
		log.Println(err)
	}
	BOT.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
	BOT.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, res.MessageID))
}

func writeMessageInDB(contextID int, messageText string) {
	err := db.SetNewMessage(contextID, messageText)
	if nil != err {
		log.Printf("Ошибка записи нового сообщения в БД, context id: '%v': %v", contextID, err)
		return
	}

	// todo: Проверить, если сообщений больше 100, то удалить самые старые.
}

func createNewContextAndSelect(contextName string, message *tgbotapi.Message) {
	response := db.SetNewContext(contextName, message.From.ID)
	// todo добавить сообщения об ответе от функции выше
	quickSendTextMessage(message.Chat.ID, response)
}

func renameContext(message *tgbotapi.Message, contextID *int) {
	if db.RenameContextByID(contextID, &message.Text) {
		// todo Исправить текст
		quickSendTextMessage(message.Chat.ID, "Контекст переименован!")
	}
}

func deleteThisMessage(message *tgbotapi.Message) {
	BOT.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
}

func checkNameContext(message *tgbotapi.Message) bool {
	if len([]rune(message.Text)) < 2 {
		// todo Изменить текст
		quickSendTextMessage(message.Chat.ID, "Имя контекста не может быть короче двух символов")
		return true
	}
	return false
}

func addWaitingUser(userID *int64, byPassword bool, userName, key *string) bool {
	return db.AddWaitingUser(userID, byPassword, userName, key)
}

func checkNewUserInWaiting(userName *string) (string, bool) {
	fstr := fmt.Sprintf("%v %v", *userName, strings.Repeat(" ", 199-len([]rune(*userName))))
	name := db.GetWaitingUser(&fstr, false)
	return name, len(name) > 0
}
