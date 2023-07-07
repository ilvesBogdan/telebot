package ai

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func Request(ctx context.Context, chanelResponse chan string, message string) {
	ctx, cancel := context.WithTimeout(ctx, time.Hour)

	req := gogpt.CompletionRequest{
		Model:     MODEL,
		MaxTokens: 2048,
		Prompt:    message,
		Stream:    true,
	}

	stream, err := CLIENT.CreateCompletionStream(ctx, req)
	if err != nil {
		log.Println("Ошибка в подключении к OpenAI: ", err)
		// todo Исправить текст сообщения
		chanelResponse <- "\n\nОшибка подключения к серверу OpenAI."
		close(chanelResponse)
		cancel()
		return
	}
	defer stream.Close()

	var response gogpt.CompletionResponse

	for {
		select {
		case <-ctx.Done():
			break
		}

		response, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Printf("Ошибка стрима OpenAI: %v\n", err)
			// todo исправить текст ошибки
			chanelResponse <- "\n\nerror"
			break
		}

		for _, part := range response.Choices {
			chanelResponse <- part.Text
		}
	}

	close(chanelResponse)
	cancel()
}
