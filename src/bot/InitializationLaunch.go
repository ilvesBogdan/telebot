package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var UPDATES_CHANEL tgbotapi.UpdatesChannel
var BOT *tgbotapi.BotAPI

/*
Init - функция инициализации бота.

Параметры:

	token - токен бота;
	timeout - время ожидания обновления;
	debug - режим отладки.

Функция создает новый экземпляр API-клиента, используя указанный токен.
Для установки режима отладки, debug-флагу присвоется true-value.
Затем, bot.GetUpdatesChan()-метод API-клиента, создаст channel, который будет слушать updates.
*/
func Init(token string, timeout int, debug bool) {
	var err error
	BOT, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Неудалось инициализировать Telegram Bot Api: ", err)
	}

	BOT.Debug = debug

	log.Printf("Authorized on account %s", BOT.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	UPDATES_CHANEL = BOT.GetUpdatesChan(u)
}

/*
Run - функция обработки входящих сообщений.

Функция прослушивает channel UPDATES_CHANEL, который создается в Init()-методе.
При получении update, сообщения или callback-query, соответствующая функция handlerIncomingMessages() или callbackInlineButton() будет запущена.
*/
func Run() {
	for update := range UPDATES_CHANEL {
		if nil != update.Message {
			go handlerIncomingMessages(update)
		} else if nil != update.CallbackQuery {
			go callbackInlineButton(update)
		}
	}
}
