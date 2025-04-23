import {
  Controller,
  Get,
  Post,
  Body,
  Query,
  Req,
  Patch,
  Param,
} from '@nestjs/common';
import { ChatService } from './chat.service';
import { Request } from 'express';

@Controller('chats')
export class ChatController {
  constructor(private readonly chatService: ChatService) {}

  @Post('private')
  createPrivateChat(@Body() body: { user1Id: number; user2Id: number }) {
    return this.chatService.createPrivateChat(body.user1Id, body.user2Id);
  }

  @Post('group')
  createGroupChat(@Body() body: { title: string; participantIds: number[] }) {
    return this.chatService.createGroupChat(body.title, body.participantIds);
  }

  @Get()
  getUserChats(@Query('userId') userId: string) {
    return this.chatService.getUserChats(Number(userId));
  }

  @Post('ai/bind')
  bindAiToChat(
    @Req() req: Request,
    @Body()
    body: {
      chatId: number;
      mode?: 'ADVISOR' | 'AUTO_REPLY';
      customPrompt?: string;
    },
  ) {
    return this.chatService.bindAiForUser(
      req,
      body.chatId,
      body.mode,
      body.customPrompt,
    );
  }

  @Patch('ai/bind/:chatId')
  updateAiBinding(
    @Req() req: Request,
    @Param('chatId') chatId: string,
    @Body()
    body: {
      mode?: 'ADVISOR' | 'AUTO_REPLY';
      customPrompt?: string;
    },
  ) {
    return this.chatService.updateAiBinding(req, Number(chatId), body);
  }
}
