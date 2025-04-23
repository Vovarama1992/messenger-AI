import { NestFactory } from '@nestjs/core';
import { WsModule } from './ws/ws.module';
import { ConfigService } from '@nestjs/config';

async function bootstrap() {
  const app = await NestFactory.create(WsModule);
  const config = app.get(ConfigService);
  const port = config.get<number>('WS_PORT_AI_AUTO') || 4003;
  await app.listen(port);
  console.log(`🚀 WS-Gateway-AUTO_REPLY listening on port ${port}`);
}
bootstrap();
