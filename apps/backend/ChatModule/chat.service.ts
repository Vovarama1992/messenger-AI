import { Injectable } from '@nestjs/common';
import { PrismaService } from 'libs/PrismaModule/prisma.service';
import { ChatType } from 'types';
import { UsersService } from '../UserModule/users.service';
import { Request } from 'express';

@Injectable()
export class ChatService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly usersService: UsersService,
  ) {}

  createPrivateChat(user1Id: number, user2Id: number) {
    return this.prisma.chat.create({
      data: {
        type: ChatType.PRIVATE,
        participants: {
          create: [{ userId: user1Id }, { userId: user2Id }],
        },
      },
      include: { participants: true },
    });
  }

  createGroupChat(title: string, participantIds: number[]) {
    return this.prisma.chat.create({
      data: {
        type: ChatType.GROUP,
        title,
        participants: {
          create: participantIds.map((userId) => ({ userId })),
        },
      },
      include: { participants: true },
    });
  }

  async bindAiForUser(
    req: Request,
    chatId: number,
    mode: 'ADVISOR' | 'AUTO_REPLY' = 'ADVISOR',
    customPrompt?: string,
  ) {
    const user = await this.usersService.authenticate(req);

    return this.prisma.aiChatBinding.upsert({
      where: {
        userId_chatId: { userId: user.id, chatId },
      },
      update: {
        enabled: true,
        mode,
        customPrompt,
      },
      create: {
        userId: user.id,
        chatId,
        enabled: true,
        mode,
        customPrompt,
      },
    });
  }

  async updateAiBinding(
    req: Request,
    chatId: number,
    data: {
      mode?: 'ADVISOR' | 'AUTO_REPLY';
      customPrompt?: string;
    },
  ) {
    const user = await this.usersService.authenticate(req);

    return this.prisma.aiChatBinding.update({
      where: {
        userId_chatId: {
          userId: user.id,
          chatId,
        },
      },
      data,
    });
  }

  getUserChats(userId: number) {
    return this.prisma.chat.findMany({
      where: {
        participants: {
          some: { userId },
        },
      },
      include: {
        participants: true,
      },
    });
  }
}
