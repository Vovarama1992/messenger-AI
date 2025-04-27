package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Vovarama1992/go-ai-service/db"
	"github.com/Vovarama1992/go-ai-service/gpt"
	"github.com/Vovarama1992/go-ai-service/pkg/types"

	ch "github.com/Vovarama1992/go-ai-service/internal/kafka/channels"

	kafkago "github.com/segmentio/kafka-go"
)

func StartAdviceWorkers(ctx context.Context, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runAdviceConsumer(ctx, id)
		}(i)
	}
}

func runAdviceConsumer(ctx context.Context, id int) {
	topic := "chat.message.ai.advice-request"
	groupId := "ai-advisor-group"
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupId,
	})

	defer reader.Close()

	fmt.Printf("👷 [Advice %d] слушает %s...\n", id, topic)
	postgres := db.NewPostgresService(db.DB)

	for {
		select {
		case <-ctx.Done():
			log.Printf("⛔ [Advice %d] завершение по сигналу", id)
			return
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("❌ [Advice %d] Kafka read error: %v", id, err)
				continue
			}

			var msg types.AiAdviceRequest
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("❌ [Advice %d] JSON error: %v", id, err)
				continue
			}

			threadId, err := postgres.EnsureThreadId(msg.TargetUserId, msg.ChatId)
			if err != nil {
				log.Printf("❌ [Advice %d] threadId error: %v", id, err)
				continue
			}

			userName, err := postgres.GetUserName(msg.SenderId)
			if err != nil {
				log.Printf("⚠️ Не удалось получить имя пользователя: %v", err)
				userName = fmt.Sprintf("userId: %d", msg.SenderId)
			}

			ch.AdviceInput <- types.EnhancedAdviceRequest{
				Request:  msg,
				ThreadId: threadId,
				UserName: userName,
			}
		}
	}
}

func StartAdviceGPTWorker(ctx context.Context, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runAdviceGPTWorker(ctx, id)
		}(i)
	}
}

func runAdviceGPTWorker(ctx context.Context, id int) {
	done := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("⛔ [Advice-GPT %d] Сигнал завершения получен", id)
			done = true

		case msg, ok := <-ch.AdviceInput:
			if !ok {
				log.Printf("🚪 [Advice-GPT %d] Канал закрыт, выходим", id)
				return
			}

			advice, err := gpt.GetAdvice(
				msg.Request.CustomPrompt,
				fmt.Sprintf("Сообщение от %s: %s", msg.UserName, msg.Request.SourceText),
				msg.ThreadId,
			)
			if err != nil {
				log.Printf("❌ [Advice-GPT %d] GPT error: %v", id, err)
				continue
			}

			ch.AdviceOutput <- types.AiAdviceResponse{
				ChatId:       msg.Request.ChatId,
				TargetUserId: msg.Request.TargetUserId,
				Advice:       advice,
				SourceText:   msg.Request.SourceText,
			}
		}

		if done {
			if len(ch.AdviceInput) == 0 {
				log.Printf("✅ [Advice-GPT %d] Все сообщения обработаны, выходим", id)
				return
			}
		}
	}
}

func StartAdviceProducerWorker(ctx context.Context, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runAdviceProducer(ctx, id)
		}(i)
	}
}

func runAdviceProducer(ctx context.Context, id int) {
	done := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("⛔ [Advice-Producer %d] Сигнал завершения получен", id)
			done = true

		case msg, ok := <-ch.AdviceOutput:
			if !ok {
				log.Printf("🚪 [Advice-Producer %d] Канал закрыт, завершаемся", id)
				return
			}

			err := SendMessage("chat.message.ai-advice", msg)
			if err != nil {
				log.Printf("❌ [Advice-Producer %d] Kafka send error: %v", id, err)
			} else {
				log.Printf("📤 [Advice-Producer %d] Отправлен совет: %+v", id, msg)
			}
		}

		if done {
			if len(ch.AdviceOutput) == 0 {
				log.Printf("✅ [Advice-Producer %d] Все сообщения обработаны, выходим", id)
				return
			}
		}
	}
}

func StartAutoreplyWorkers(ctx context.Context, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go runAutoreplyConsumer(ctx, wg, i)
	}
}

func runAutoreplyConsumer(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	topic := "chat.message.ai.autoreply-request"
	groupId := "ai-autoreply-group"
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupId,
	})

	defer reader.Close()

	fmt.Printf("🤖 [AutoReply %d] слушает %s...\n", id, topic)

	postgres := db.NewPostgresService(db.DB)

	for {
		select {
		case <-ctx.Done():
			log.Printf("⛔ [AutoReply %d] Завершение по сигналу", id)
			return

		default:
			m, err := reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Printf("⛔ [AutoReply %d] FetchMessage отменен", id)
					return
				}
				log.Printf("❌ [AutoReply %d] Kafka error: %v", id, err)
				continue
			}

			var msg types.AiAutoreplyRequest
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("❌ [AutoReply %d] JSON error: %v", id, err)
				continue
			}

			threadId, err := postgres.EnsureThreadId(msg.TargetUserId, msg.ChatId)
			if err != nil {
				log.Printf("❌ [AutoReply %d] threadId error: %v", id, err)
				continue
			}

			userName, err := postgres.GetUserName(msg.SenderId)
			if err != nil {
				log.Printf("⚠️ [AutoReply %d] Не удалось получить имя пользователя: %v", id, err)
				userName = fmt.Sprintf("userId: %d", msg.SenderId)
			}

			select {
			case <-ctx.Done():
				log.Printf("⛔ [AutoReply %d] Завершение перед записью в канал", id)
				return
			case ch.AutoReplyInput <- types.EnhancedAutoreplyRequest{
				Request:  msg,
				ThreadId: threadId,
				UserName: userName,
			}:

			}
		}
	}
}

func StartAutoreplyGPTWorkers(ctx context.Context, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go runAutoreplyGPTWorker(ctx, wg, i)
	}
}

func runAutoreplyGPTWorker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	done := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("⛔ [AutoReply-GPT %d] Сигнал завершения получен", id)
			done = true

		case msg, ok := <-ch.AutoReplyInput:
			if !ok {
				log.Printf("🚪 [AutoReply-GPT %d] Канал закрыт, выходим", id)
				return
			}

			formatted := fmt.Sprintf("Сообщение от %s: %s", msg.UserName, msg.Request.Text)

			reply, err := gpt.GetAutoreply(msg.Request.CustomPrompt, formatted, msg.ThreadId)
			if err != nil {
				log.Printf("❌ GPT error (AutoReply-GPT worker %d): %v", id, err)
				continue
			}

			select {
			case <-ctx.Done():
				log.Printf("⛔ [AutoReply-GPT %d] Завершение перед записью в канал", id)
				return
			case ch.AutoReplyOutput <- types.AiAutoreplyResponse{
				ChatId:       msg.Request.ChatId,
				TargetUserId: msg.Request.TargetUserId,
				Text:         reply,
				SenderId:     msg.Request.TargetUserId,
			}:

			}
		}

		if done && len(ch.AutoReplyInput) == 0 {
			log.Printf("✅ [AutoReply-GPT %d] Все сообщения обработаны, выходим", id)
			return
		}
	}
}

func StartAutoReplySenderWorkers(ctx context.Context, wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		go runAutoReplySender(ctx, wg, i)
	}
}

func runAutoReplySender(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	done := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("⛔ [AutoReply-Sender %d] Сигнал завершения получен", id)
			done = true

		case msg, ok := <-ch.AutoReplyOutput:
			if !ok {
				log.Printf("🚪 [AutoReply-Sender %d] Канал закрыт, выходим", id)
				return
			}

			SendAutoreply(msg)
			SendPersist(msg.ChatId, msg.SenderId, msg.Text)
		}

		if done && len(ch.AutoReplyOutput) == 0 {
			log.Printf("✅ [AutoReply-Sender %d] Все сообщения обработаны, выходим", id)
			return
		}
	}
}
func SendPersist(chatId int, senderId int, text string) {
	_ = SendMessage("chat.message.persist", map[string]interface{}{
		"chatId":   chatId,
		"senderId": senderId,
		"text":     text,
	})
}
