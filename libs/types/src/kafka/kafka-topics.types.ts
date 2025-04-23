import { CreateMessageDto } from '../message/create-message.dto';

export type KafkaTopics = {
  'chat.message.persist': CreateMessageDto;
  'chat.message.forward': CreateMessageDto & { senderSocketId: string };

  'chat.message.ai.advice-request': {
    chatId: number;
    targetUserId: number;
    sourceText: string;
    customPrompt?: string;
  };

  'chat.message.ai-advice': {
    chatId: number;
    targetUserId: number;
    advice: string;
    sourceText: string;
  };

  'chat.message.ai.autoreply-request': CreateMessageDto & {
    targetUserId: number;
    customPrompt?: string;
  };
};
