Мессенджер с микросервисной архитектурой подключенный к AI:

NestJS (REST API + WebSocket Gateway)

Kafka (балансировка, очереди)

Go-сервис (GPT-обработчик, многопоточность)

PostgreSQL + Prisma

Docker / Docker Compose

📚 Архитектура
backend — основной API

ws-gateway-1 и ws-gateway-2 — WebSocket-сервисы (нагрузка делится)

websocket-ai-auto — слушает chat.message.ai.autoreply-request, шлёт автоответ

websocket-ai-advice — слушает chat.message.ai.advice-request, возвращает совет

go-ai-service — Go-сервис, общается с OpenAI GPT, шлёт ответ обратно в Kafka

Kafka брокер — точка соединения всех

🔌 Поток сообщений
Пользователь шлёт сообщение через WebSocket.

Kafka пересылает другим WebSocket'ам и Go-сервису.

Go-сервис обращается к GPT.

Результат идёт обратно через Kafka → нужный WebSocket → клиенту.

🧪 Технологии
NestJS, WebSocket, Prisma

Go + goroutines + GPT API

Kafka 

PostgreSQL

Docker (всё через docker-compose)

TypeScript, tsconfig-paths, монорепа

🚀 Как запустить
git clone ...
cd myGoogle
docker-compose up --build
