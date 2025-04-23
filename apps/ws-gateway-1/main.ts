import { NestFactory } from '@nestjs/core';
import { WsModule } from './ws/ws.module';
import { ConfigService } from '@nestjs/config';

async function bootstrap() {
  const app = await NestFactory.create(WsModule);
  const config = app.get(ConfigService);
  const port = config.get<number>('WS_PORT_1') || 4001;
  await app.listen(port);
  console.log(`🚀 WS-Gateway-1 listening on port ${port}`);
}
bootstrap();
