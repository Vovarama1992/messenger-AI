import {
  WebSocketGateway,
  WebSocketServer,
  OnGatewayInit,
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { KafkaService } from 'libs/KafkaModule/kafka.service';
import { Logger } from '@nestjs/common';

@WebSocketGateway({ cors: { origin: '*' } })
export class WebSocketAdviceGateway implements OnGatewayInit {
  @WebSocketServer()
  server: Server;

  private readonly logger = new Logger(WebSocketAdviceGateway.name);

  constructor(private readonly kafkaService: KafkaService) {}

  async afterInit() {
    const consumer = this.kafkaService.getConsumer();
    await consumer.connect();
    await consumer.subscribe({ topic: 'chat.message.ai-advice' });

    await consumer.run({
      eachMessage: async ({ message }) => {
        const data = JSON.parse(message.value.toString());

        this.logger.log(`📥 Kafka → AI Advice: ${JSON.stringify(data)}`);

        this.server.to(`user:${data.targetUserId}`).emit('aiAdvice', {
          chatId: data.chatId,
          text: data.advice,
          sourceText: data.sourceText,
        });
      },
    });
  }

  handleConnection(client: Socket) {
    const userId = client.handshake.query.userId;
    if (userId) {
      client.join(`user:${userId}`);
      this.logger.log(`✅ Client ${client.id} joined user:${userId}`);
    } else {
      this.logger.warn(`⚠️ Client ${client.id} connected without userId`);
    }
  }

  handleDisconnect(client: Socket) {
    this.logger.log(`❌ Client disconnected: ${client.id}`);
  }
}
