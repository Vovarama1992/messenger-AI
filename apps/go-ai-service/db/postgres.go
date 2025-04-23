package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/sashabaranov/go-openai"
)

var DB *pgxpool.Pool

func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	DB = pool
	return nil
}

type PostgresService struct {
	pool   *pgxpool.Pool
	client *openai.Client
}

func NewPostgresService(pool *pgxpool.Pool) *PostgresService {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	return &PostgresService{
		pool:   pool,
		client: client,
	}
}

func (p *PostgresService) EnsureThreadId(userId int, chatId int) (string, error) {
	ctx := context.Background()

	var threadId string
	err := p.pool.QueryRow(ctx, `
		SELECT "threadId" FROM "AiChatBinding"
		WHERE "userId" = $1 AND "chatId" = $2
	`, userId, chatId).Scan(&threadId)

	if err == nil && threadId != "" {
		return threadId, nil
	}

	thread, err := p.client.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	_, err = p.pool.Exec(ctx, `
		UPDATE "AiChatBinding"
		SET "threadId" = $1
		WHERE "userId" = $2 AND "chatId" = $3
	`, thread.ID, userId, chatId)

	if err != nil {
		return "", fmt.Errorf("failed to update threadId: %w", err)
	}

	return thread.ID, nil
}
