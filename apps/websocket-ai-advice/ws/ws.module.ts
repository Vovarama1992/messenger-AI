import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { KafkaModule } from 'libs/KafkaModule/kafka.module';
import { PrismaModule } from 'libs/PrismaModule/prisma.module';
import { WebSocketAdviceGateway } from './ws.gateway';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    KafkaModule,
    PrismaModule,
  ],
  providers: [WebSocketAdviceGateway],
})
export class WsModule {}
