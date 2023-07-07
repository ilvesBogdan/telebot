package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/ai"
)

func sendingMessagesFromAI(contextText, messageText string, message tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Hour)

	responseFromAI := make(chan string)

	// Составление сообщения для ИИ
	var contextAndText string
	if len(contextText) > 4 {
		contextAndText = fmt.Sprintf("%v (%v)", contextText, messageText)
	} else {
		contextAndText = messageText
	}

	// Отправляем запрос к ИИ
	go ai.Request(ctx, responseFromAI, contextAndText)

	// todo Добавить сообщения загрузки
	msgResp := quickSendTextMessage(message.Chat.ID, "Loading ...")

	var delay = int64(400)               // задержка перед отправкой в телеграм в милисекундах
	var lenghtMsgText = int(7)           // длина полученого от ИИ сообщения на предыдущей итерации
	var msgTimeOut = time.Now()          // Текущее время
	var msgText = strings.Builder{}      // Бефер для частей сообщения полученого от ИИ
	var responseMsg tgbotapi.APIResponse // Ответ от сервера телеграм

	for text := range responseFromAI {
		msgText.WriteString(text)

		select {
		case <-ctx.Done():
			log.Println("Таймаут выполнения рутины.")
			// todo Исправить текст сообщения
			msgText.WriteString("\n\nТаймаут генерации.")
			break
		}

		if time.Since(msgTimeOut).Milliseconds() > delay && msgText.Len() > lenghtMsgText {
			responseMsg = editAiMessage(message.Chat, &msgResp, fmt.Sprintf("%v\n\n...", msgText.String()))

			msgTimeOut = time.Now()
			if responseMsg.Ok {
				delay = delay + int64(10)
				lenghtMsgText = msgText.Len() + 4
				continue
			}

			switch responseMsg.ErrorCode {
			case 1:
				cancel()
				return
			case 500:
				log.Println("Ошибка сервера телеграм: ", responseMsg)
				cancel()
				return
			case 429:
				// Если в телеграм было отправлено слишком много запросов
				delay = int64(responseMsg.Parameters.RetryAfter * 1000)
			default:
				if responseMsg.ErrorCode >= 420 && responseMsg.ErrorCode < 500 {
					delay = int64(responseMsg.Parameters.RetryAfter * 1001)
					continue
				}
				log.Println("Неизвестная ошибка телеграм: ", responseMsg)
				cancel()
				return
			}
		}
	}

	if !responseMsg.Ok && responseMsg.ErrorCode >= 420 && responseMsg.ErrorCode < 500 {
		delay = int64(responseMsg.Parameters.RetryAfter*1000 + 100)
	} else {
		delay = int64(500)
	}

	var attempts = 0 // Количество попыток редактирования сообщения

	for {
		if time.Since(msgTimeOut).Milliseconds() > delay {
			responseMsg = editAiMessage(message.Chat, &msgResp, msgText.String())
			if responseMsg.Ok {
				cancel()
				return
			}

			// Если в телеграм было отправлено слишком много запросов
			if responseMsg.ErrorCode >= 420 && responseMsg.ErrorCode < 500 {
				delay = int64(responseMsg.Parameters.RetryAfter * 1001)
			} else {
				delay = int64(1000)
			}

			msgTimeOut = time.Now()
			attempts = attempts + 1
		}

		if attempts > 20 {
			log.Println("Таймаут отправки сообщения от ИИ: ", responseMsg)
			cancel()
			return
		}
	}
}

func editAiMessage(chat *tgbotapi.Chat, msgResp *tgbotapi.Message, text string) tgbotapi.APIResponse {
	response, err := BOT.Request(tgbotapi.NewEditMessageText(chat.ID, msgResp.MessageID, text))
	if nil != err {
		log.Println("Ошибка редактирования сообщения в телеграм, от ИИ: ", err)
		return tgbotapi.APIResponse{Ok: false, ErrorCode: 1}
	}
	return *response
}
