import { Module } from '@nestjs/common';
import { ChatService } from './chat.service';
import { ChatController } from './chat.controller';
import { PrismaModule } from 'libs/PrismaModule/prisma.module';

@Module({
  imports: [PrismaModule],
  providers: [ChatService],
  controllers: [ChatController],
})
export class ChatModule {}
