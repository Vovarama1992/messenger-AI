import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { WebSocketGateway2 } from './ws.gateway';
import { KafkaModule } from 'libs/KafkaModule/kafka.module';
import { PrismaModule } from 'libs/PrismaModule/prisma.module';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    KafkaModule,
    PrismaModule,
  ],
  providers: [WebSocketGateway2],
})
export class WsModule {}
