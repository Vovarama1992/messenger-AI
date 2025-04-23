import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { WebSocketAIAutoReplyGateway } from './websocket-ai-auto.gateway';
import { KafkaModule } from 'libs/KafkaModule/kafka.module';

@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true }), KafkaModule],
  providers: [WebSocketAIAutoReplyGateway],
})
export class WsModule {}
