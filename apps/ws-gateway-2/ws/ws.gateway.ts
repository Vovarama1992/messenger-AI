import {
  WebSocketGateway,
  WebSocketServer,
  SubscribeMessage,
  MessageBody,
  ConnectedSocket,
} from '@nestjs/websockets';
import { OnModuleInit } from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { KafkaService } from 'libs/KafkaModule/kafka.service';
import { Logger } from '@nestjs/common';
import { PrismaService } from 'libs/PrismaModule/prisma.service';

@WebSocketGateway({ cors: { origin: '*' } })
export class WebSocketGateway1 implements OnModuleInit {
  @WebSocketServer()
  server: Server;

  private readonly logger = new Logger(WebSocketGateway1.name);

  constructor(
    private readonly kafkaService: KafkaService,
    private readonly prisma: PrismaService,
  ) {}

  async onModuleInit() {
    const consumer = this.kafkaService.getConsumer();
    await consumer.connect();
    await consumer.subscribe({ topic: 'chat.message.forward' });

    await consumer.run({
      eachMessage: async ({ message }) => {
        const data = JSON.parse(message.value.toString());
        this.logger.log(`📥 Kafka → WS1: ${JSON.stringify(data)}`);
        if (this.server.sockets.sockets.has(data.senderSocketId)) {
          return;
        }
        this.broadcastMessage(data);
      },
    });
  }

  handleConnection(client: Socket) {
    this.logger.log(`✅ Client connected: ${client.id}`);
  }

  handleDisconnect(client: Socket) {
    this.logger.log(`❌ Client disconnected: ${client.id}`);
  }

  @SubscribeMessage('sendMessage')
  async handleMessage(
    @MessageBody() data: { chatId: number; text: string; senderId: number },
    @ConnectedSocket() client: Socket,
  ) {
    const enriched = {
      ...data,
      senderSocketId: client.id,
    };
    this.logger.log(`📨 From ${client.id}: ${JSON.stringify(data)}`);
    await this.kafkaService.sendMessage('chat.message.persist', data);

    await this.kafkaService.sendMessage('chat.message.forward', enriched);

    this.broadcastMessage(enriched);

    const aiBindings = await this.prisma.aiChatBinding.findMany({
      where: { chatId: Number(data.chatId), enabled: true },
    });

    for (const binding of aiBindings) {
      if (binding.mode === 'ADVISOR') {
        await this.kafkaService.sendMessage('chat.message.ai.advice-request', {
          chatId: data.chatId,
          sourceText: data.text,
          targetUserId: binding.userId,
          customPrompt: binding.customPrompt || undefined,
        });
      } else if (binding.mode === 'AUTO_REPLY') {
        await this.kafkaService.sendMessage(
          'chat.message.ai.autoreply-request',
          {
            ...data,
            targetUserId: binding.userId,
            customPrompt: binding.customPrompt || undefined,
          },
        );
      }
    }
  }

  async broadcastMessage(message: { chatId: number; text: string }) {
    this.server.to(String(message.chatId)).emit('newMessage', message);
  }
}
