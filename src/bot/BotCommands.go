package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commands(text *LanguagePackageOfMessages, update *tgbotapi.Update) (isCommand bool) {

	// Нажата ли клавиша выбора контекста
	if len(update.Message.Text) > 6 {
		switch update.Message.Text[:4] {
		// Кейс с пробелом, т.к. размер в байтах у эмоджи отличается
		case "✅ ":
			// Нажат текущий контекст
			isCommand = clickSelectedContext(text, update.Message)
			return
		case "💬":
			// Нажат один из существующих контекстов
			isCommand = clickContext(text, update.Message)
			return
		}
	}

	isCommand = update.Message.IsCommand()

	if !isCommand {
		return
	}

	switch update.Message.Command() {

	// Вывести сообщение при первом запуске бота
	case "start":
		// todo Реализовать

	// Режим диалога без контекста
	case "without_context":
		commandWithoutContext(update.Message)

	// Создаем новый контекст
	case "new":
		commandNewContext(text, update.Message)

	// Показываем список контекстов
	case "show_all_contexts":
		showMenuContexts(update.Message)

	// Скрываем список контекстов
	case "hide_contexts":
		hideMenuContexts(update.Message)

	// Конманды для администратора бота
	case "adduser":
		commandAddUser(update.Message)

	case "showreg":
		commandShowReg(update.Message)

	case "rmuser":
		// todo Реализовать удаление пользователей
		commandRmUser(update)

	case "showmsgs":
		// todo Реализовать просмотр сообщений
		commandShowMsgs(update)
	}

	return isCommand
}

func commandNewContext(text *LanguagePackageOfMessages, message *tgbotapi.Message) {
	contextName := message.CommandArguments()

	if contextName == "" {
		setExpectedMessageTypeForUser(2, &message.From.ID)
		quickSendTextMessage(message.Chat.ID, text.msg("EnterContextName"))
		return
	} else if len([]rune(contextName)) < 2 {
		// todo Изменить текст
		quickSendTextMessage(message.Chat.ID, "Имя контекста не может быть короче двух символов")
		return
	}

	createNewContextAndSelect(contextName, message)
}

func commandWithoutContext(message *tgbotapi.Message) {
	// todo Исправить текст
	if withoutContext(&message.From.ID) {
		quickSendTextMessage(message.Chat.ID, "Теперь без контекста")
	} else {
		quickSendTextMessage(message.Chat.ID, "Ошибка")
	}
	deleteThisMessage(message)
}

func commandAddUser(message *tgbotapi.Message) {

	if !userIsAdmin(message.From) {
		return
	}

	arg := strings.Split(message.CommandArguments(), "@")

	if len(arg) != 2 {
		quickSendTextMessage(message.Chat.ID, "error: nick@username")
		return
	}

	nick, username := arg[0], arg[1]

	addWaitingUser(&message.From.ID, false, &nick, &username)
	quickSendTextMessage(message.Chat.ID, fmt.Sprintf("😴 %v", nick))
	BOT.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
}

func commandShowReg(message *tgbotapi.Message) {

	if !userIsAdmin(message.From) {
		return
	}

	// todo Реализовать

	// if len(usersWaitingToRegister) < 1 {
	// 	quickSendTextMessage(message.Chat.ID, "null")
	// 	return
	// }

	// var buttons [][]tgbotapi.InlineKeyboardButton
	// for username, nick := range usersWaitingToRegister {
	// 	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
	// 		tgbotapi.NewInlineKeyboardButtonData(
	// 			fmt.Sprintf("❌ %v", nick),
	// 			fmt.Sprintf("rmRegUser:%v", username),
	// 		),
	// 	))
	// }
	// buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
	// 	tgbotapi.NewInlineKeyboardButtonData("❌", "close"),
	// ))

	// msg := tgbotapi.NewMessage(message.Chat.ID, "...")
	// msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	// if _, err := BOT.Send(msg); err != nil {
	// 	log.Panic(err)
	// }
	// BOT.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
}

func commandRmUser(update *tgbotapi.Update) {

	if !userIsAdmin(update.Message.From) {
		return
	}
}

func commandShowMsgs(update *tgbotapi.Update) {

	if !userIsAdmin(update.Message.From) {
		return
	}
}
