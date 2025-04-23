import { Injectable } from '@nestjs/common';
import { PrismaService } from 'libs/PrismaModule/prisma.service';
import { CreateMessageDto } from 'types';

@Injectable()
export class MessageService {
  constructor(private readonly prisma: PrismaService) {}

  createMessage(dto: CreateMessageDto) {
    return this.prisma.message.create({
      data: {
        chatId: Number(dto.chatId),
        senderId: dto.senderId,
        text: dto.text,
      },
    });
  }

  getMessagesByChat(chatId: number, limit = 50, offset = 0) {
    return this.prisma.message.findMany({
      where: { chatId },
      orderBy: { createdAt: 'asc' },
      take: limit,
      skip: offset,
    });
  }
}
