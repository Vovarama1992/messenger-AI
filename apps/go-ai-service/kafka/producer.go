package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func SendAdvice(msg AiAdviceResponse) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "chat.message.ai-advice",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("❌ Kafka marshal error: %v", err)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: bytes,
	})
	if err != nil {
		log.Printf("❌ Kafka write error: %v", err)
	} else {
		log.Printf("📤 Совет отправлен в Kafka: %s", string(bytes))
	}
}

func SendAutoreply(msg AiAutoreplyRequest) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    "chat.message.forward",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("❌ Autoreply marshal error: %v", err)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{Value: bytes})
	if err != nil {
		log.Printf("❌ Kafka write error (autoreply): %v", err)
	} else {
		log.Printf("📤 AutoReply отправлен: %s", string(bytes))
	}
}

func SendMessage(topic string, payload any) error {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Kafka marshal error: %v", err)
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{Value: bytes})
	if err != nil {
		log.Printf("❌ Kafka write error: %v", err)
		return err
	}

	log.Printf("📤 Kafka → [%s]: %s", topic, string(bytes))
	return nil
}
