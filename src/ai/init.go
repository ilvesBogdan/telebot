package ai

import (
	gogpt "github.com/sashabaranov/go-gpt3"
)

var CLIENT *gogpt.Client
var MODEL string

/*
Init - функция инициализации бота.

Параметры:

	token - токен для авторизации в GPT-3 API
	model - название модели GPT-3, которую будет использовать AI-бот

Функция создает соединение с GPT-3 API, устанавливая указанный токен,
а также сохраняя название модели GPT-3.
*/
func Init(token string, model string) {
	CLIENT = gogpt.NewClient(token)
	MODEL = model
}
