import { Module, OnModuleInit } from '@nestjs/common';
import { MessageService } from './message.service';
import { MessageController } from './message.controller';
import { KafkaService } from 'libs/KafkaModule/kafka.service';
import { PrismaModule } from 'libs/PrismaModule/prisma.module';

@Module({
  imports: [PrismaModule],
  controllers: [MessageController],
  providers: [MessageService, KafkaService],
})
export class MessageModule implements OnModuleInit {
  constructor(
    private readonly kafkaService: KafkaService,
    private readonly messageService: MessageService,
  ) {}

  async onModuleInit() {
    const consumer = this.kafkaService.getConsumer();

    await consumer.connect();
    await consumer.subscribe({
      topic: 'chat.message.persist',
      fromBeginning: false,
    });

    await consumer.run({
      eachMessage: async ({ message }) => {
        const data = JSON.parse(message.value.toString());
        await this.messageService.createMessage(data);
      },
    });
  }
}
