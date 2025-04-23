package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Vovarama1992/go-ai-service/db"
	"github.com/Vovarama1992/go-ai-service/gpt"
	"github.com/segmentio/kafka-go"
)

type AiAdviceRequest struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	SourceText   string `json:"sourceText"`
	CustomPrompt string `json:"customPrompt,omitempty"`
}

type AiAdviceResponse struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	Advice       string `json:"advice"`
	SourceText   string `json:"sourceText"`
}

type AiAutoreplyRequest struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	Text         string `json:"text"`
	SenderId     int    `json:"senderId"`
	CustomPrompt string `json:"customPrompt,omitempty"`
}

func StartAdviceWorkers(count int) {
	for i := 0; i < count; i++ {
		go startAdviceConsumer(i)
	}
}

func StartAutoreplyWorkers(count int) {
	for i := 0; i < count; i++ {
		go startAutoreplyConsumer(i)
	}
}

func startAdviceConsumer(id int) {
	topic := "chat.message.ai.advice-request"
	groupId := "ai-advisor-group"
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupId,
	})

	fmt.Printf("\U0001f477 [Advice %d] слушает %s...\n", id, topic)

	postgres := db.NewPostgresService(db.DB)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ [Advice %d] Kafka read error: %v", id, err)
			continue
		}

		var msg AiAdviceRequest
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("❌ [Advice %d] JSON error: %v", id, err)
			continue
		}

		fmt.Printf("📥 [Advice %d] input: %+v\n", id, msg)

		threadId, err := postgres.EnsureThreadId(msg.TargetUserId, msg.ChatId)
		if err != nil {
			log.Printf("❌ [Advice %d] threadId error: %v", id, err)
			continue
		}

		advice, err := gpt.GetAdvice(msg.CustomPrompt, msg.SourceText, threadId)
		if err != nil {
			log.Printf("❌ [Advice %d] GPT error: %v", id, err)
			continue
		}

		SendAdvice(AiAdviceResponse{
			ChatId:       msg.ChatId,
			TargetUserId: msg.TargetUserId,
			Advice:       advice,
			SourceText:   msg.SourceText,
		})
	}
}

func SendPersist(chatId int, senderId int, text string) {
	_ = SendMessage("chat.message.persist", map[string]interface{}{
		"chatId":   chatId,
		"senderId": senderId,
		"text":     text,
	})
}

func startAutoreplyConsumer(id int) {
	topic := "chat.message.ai.autoreply-request"
	groupId := "ai-autoreply-group"
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupId,
	})

	fmt.Printf("🤖 [AutoReply %d] слушает %s...\n", id, topic)

	postgres := db.NewPostgresService(db.DB)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ [AutoReply %d] Kafka error: %v", id, err)
			continue
		}

		var msg AiAutoreplyRequest
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("❌ [AutoReply %d] JSON error: %v", id, err)
			continue
		}

		fmt.Printf("📥 [AutoReply %d] input: %+v\n", id, msg)

		threadId, err := postgres.EnsureThreadId(msg.TargetUserId, msg.ChatId)
		if err != nil {
			log.Printf("❌ [AutoReply %d] threadId error: %v", id, err)
			continue
		}

		reply, err := gpt.GetAutoreply(msg.CustomPrompt, msg.Text, threadId)
		if err != nil {
			log.Printf("❌ [AutoReply %d] GPT error: %v", id, err)
			continue
		}

		SendAutoreply(AiAutoreplyRequest{
			ChatId:       msg.ChatId,
			TargetUserId: msg.TargetUserId,
			Text:         reply,
			SenderId:     msg.TargetUserId,
		})
		SendPersist(msg.ChatId, msg.TargetUserId, reply)
	}
}
