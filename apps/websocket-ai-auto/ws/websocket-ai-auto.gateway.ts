import { WebSocketGateway, WebSocketServer } from '@nestjs/websockets';
import { OnModuleInit } from '@nestjs/common';
import { Server } from 'socket.io';
import { KafkaService } from 'libs/KafkaModule/kafka.service';
import { Logger } from '@nestjs/common';

@WebSocketGateway({ cors: { origin: '*' } })
export class WebSocketAIAutoReplyGateway implements OnModuleInit {
  @WebSocketServer()
  server: Server;

  private readonly logger = new Logger(WebSocketAIAutoReplyGateway.name);

  constructor(private readonly kafkaService: KafkaService) {}

  async onModuleInit() {
    const consumer = this.kafkaService.getConsumer();
    await consumer.connect();
    await consumer.subscribe({ topic: 'chat.message.forward' });

    await consumer.run({
      eachMessage: async ({ message }) => {
        const data = JSON.parse(message.value.toString());

        // Проверим, что это AI-автоответ
        const isAutoReply =
          typeof data.targetUserId === 'number' &&
          data.senderId === data.targetUserId;

        if (!isAutoReply) return;

        this.server.to(String(data.chatId)).emit('newMessage', {
          chatId: data.chatId,
          senderId: data.senderId,
          text: data.text,
          autoReply: true,
        });

        this.logger.log(
          `🤖 AUTO_REPLY → Chat ${data.chatId}, user ${data.senderId}`,
        );
      },
    });
  }
}
