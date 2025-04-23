import { NestFactory } from '@nestjs/core';
import { WsModule } from './ws/ws.module';
import { ConfigService } from '@nestjs/config';

async function bootstrap() {
  const app = await NestFactory.create(WsModule);
  const config = app.get(ConfigService);
  const port = config.get<number>('WS_PORT_AI_ADVICE') || 4004;
  await app.listen(port);
  console.log(`🚀 WS-Gateway-AI-ADVICE listening on port ${port}`);
}
bootstrap();
