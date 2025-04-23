import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { WebSocketGateway1 } from './ws.gateway';
import { KafkaModule } from 'libs/KafkaModule/kafka.module';
import { PrismaModule } from 'libs/PrismaModule/prisma.module';
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    KafkaModule,
    PrismaModule,
  ],
  providers: [WebSocketGateway1],
})
export class WsModule {}
