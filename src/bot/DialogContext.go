package bot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/ai"
	"gitlab.com/_Bogdan/chat-gpt-telegram-bot/src/db"
)

var replace = regexp.MustCompile(`[^A-Za-zА-Яа-я ]`)
var lastword = regexp.MustCompile(`,?[A-Za-zА-Яа-я ]+$`)

func getContext(contextID *int) (context db.Context) {
	context = db.GetContextsByID(contextID)
	return
}

func clickSelectedContext(text *LanguagePackageOfMessages, message *tgbotapi.Message) bool {
	contexts, err := db.GetContextsByTid(message.From.ID)
	if nil != err {
		log.Panic("Неудалось получить данные контекстов из ДБ:", err)
		return false
	}
	// Обрезаем эмоджи и пробел
	clickedContextName := message.Text[4:]
	for _, ct := range contexts {
		ct.TrimSpace()
		if clickedContextName == ct.Name {
			// Сообщение конфигурации текущего контекста
			go hideMenuContexts(message)
			sendConfigContextMessage(message, &ct)
			return true
		}
	}
	return false
}

func clickContext(text *LanguagePackageOfMessages, message *tgbotapi.Message) bool {
	contexts, err := db.GetContextsByTid(message.From.ID)
	if nil != err {
		log.Panic("Неудалось получить данные контекстов из ДБ:", err)
		return false
	}
	// Обрезаем эмоджи и пробел
	clickedContextName := message.Text[5:]
	for _, ct := range contexts {
		ct.TrimSpace()
		if clickedContextName == ct.Name {
			// Обновление контекста
			var msg tgbotapi.MessageConfig
			if db.SetUserContext(ct.Id, message.From.ID) {
				msg = tgbotapi.NewMessage(message.Chat.ID, text.msg("ContextChanged"))
			} else {
				msg = tgbotapi.NewMessage(message.Chat.ID, text.msg("error"))
			}
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := BOT.Send(msg); err != nil {
				log.Panic(err)
			}
			return true
		}
	}
	return false
}

func withoutContext(UserID *int64) bool {
	return db.SetUserContext(0, *UserID)
}

func setNewContext(oldContextText, message string, contextID int) {

	// Не обрабатываем короткие сооющения
	if len([]rune(message)) < 10 {
		return
	}

	newWordsClice := getWords(keyWordsFromAI(message))
	oldWordsClice := getWords(oldContextText)
	differenClice := difference(newWordsClice, oldWordsClice)

	if len(differenClice) < 1 {
		return
	}

	text := strings.Join(append(differenClice, oldWordsClice...), " ")

	for len([]rune(text)) > 255 {
		text = lastword.ReplaceAllString(text, "")
	}

	if len([]rune(text)) > 0 {
		db.SetNewContextText(contextID, text)
	}
}

func keyWordsFromAI(text string) string {
	responseFromAI := make(chan string)
	message := fmt.Sprintf("Выведи ключевые слова из строки \"%v\" в порядке важности.", text)
	go ai.Request(context.Background(), responseFromAI, message)
	var msgText = strings.Builder{}
	for text := range responseFromAI {
		msgText.WriteString(text)
	}
	return msgText.String()
}

func difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func getWords(text string) []string {
	var result []string
	text = strings.ToLower(text)
	text = replace.ReplaceAllString(text, "")
	for _, w := range strings.Split(text, " ") {
		if len([]rune(w)) > 3 {
			result = append(result, strings.TrimSpace(w))
		}
	}
	return result
}
