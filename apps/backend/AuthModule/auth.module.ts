import { Module } from '@nestjs/common';
import { AuthService } from './auth.service';
import { PrismaModule } from 'libs/PrismaModule/prisma.module';
import { JwtModule } from 'libs/JwtModule/jwt.module';
import { AuthController } from './auth.controller';

@Module({
  imports: [PrismaModule, JwtModule],
  providers: [AuthService],
  controllers: [AuthController],
  exports: [AuthService],
})
export class AuthModule {}
