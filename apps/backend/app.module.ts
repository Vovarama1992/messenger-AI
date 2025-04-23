import { Module } from '@nestjs/common';
import { UsersModule } from './UserModule/users.module';
import { JwtModule } from '../../libs/JwtModule/jwt.module';
import { ConfigModule } from '@nestjs/config';
import { PrismaModule } from '../../libs/PrismaModule/prisma.module';
import { AuthModule } from './AuthModule/auth.module';
import { KafkaModule } from 'libs/KafkaModule/kafka.module';
import { MessageModule } from './MessageModule/message.module';
import { ChatModule } from './ChatModule/chat.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),

    AuthModule,

    JwtModule,
    KafkaModule,
    MessageModule,
    ChatModule,

    UsersModule,

    PrismaModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
