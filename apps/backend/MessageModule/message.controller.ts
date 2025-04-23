import { Controller, Get, Post, Body, Query } from '@nestjs/common';
import { MessageService } from './message.service';
import { CreateMessageDto } from 'types';

@Controller('messages')
export class MessageController {
  constructor(private readonly messageService: MessageService) {}

  @Post()
  create(@Body() dto: CreateMessageDto) {
    return this.messageService.createMessage(dto);
  }

  @Get()
  getChatMessages(
    @Query('chatId') chatId: string,
    @Query('limit') limit?: string,
    @Query('offset') offset?: string,
  ) {
    return this.messageService.getMessagesByChat(
      Number(chatId),
      Number(limit) || 50,
      Number(offset) || 0,
    );
  }
}
