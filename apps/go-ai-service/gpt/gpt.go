package gpt

import (
	"context"
	"fmt"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var (
	client      = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	assistantID = os.Getenv("OPENAI_ASSISTANT_ID")
)

func GetAdvice(prompt string, contextText string, threadId string) (string, error) {
	ctx := context.Background()

	_, err := client.CreateMessage(ctx, threadId, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt + "\n\nПользователь написал: " + contextText,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при добавлении сообщения: %w", err)
	}

	run, err := client.CreateRun(ctx, threadId, openai.RunRequest{
		AssistantID: assistantID,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при запуске run: %w", err)
	}

	for {
		r, _ := client.RetrieveRun(ctx, threadId, run.ID)
		if r.Status == openai.RunStatusCompleted {
			break
		}
		time.Sleep(1 * time.Second)
	}

	messages, err := client.ListThreadMessages(ctx, threadId, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении сообщений: %w", err)
	}

	if len(messages.Messages) == 0 || len(messages.Messages[0].Content) == 0 {
		return "", fmt.Errorf("ответ не найден")
	}

	return messages.Messages[0].Content[0].Text.Value, nil
}

func GetAutoreply(prompt string, contextText string, threadId string) (string, error) {
	ctx := context.Background()

	_, err := client.CreateMessage(ctx, threadId, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt + "\n\nОтветь на сообщение: " + contextText,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при добавлении сообщения: %w", err)
	}

	run, err := client.CreateRun(ctx, threadId, openai.RunRequest{
		AssistantID: assistantID,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при запуске run: %w", err)
	}

	for {
		r, _ := client.RetrieveRun(ctx, threadId, run.ID)
		if r.Status == openai.RunStatusCompleted {
			break
		}
		time.Sleep(1 * time.Second)
	}

	messages, err := client.ListMessages(ctx, threadId, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении сообщений: %w", err)
	}

	if len(messages.Messages) == 0 || len(messages.Messages[0].Content) == 0 {
		return "", fmt.Errorf("ответ не найден")
	}

	return messages.Messages[0].Content[0].Text.Value, nil
}
