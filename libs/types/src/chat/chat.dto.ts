import { ApiProperty } from '@nestjs/swagger';
import { ChatType } from './chat-type.enum';

export class ChatDto {
  @ApiProperty({ example: 1 })
  id: number;

  @ApiProperty({ enum: ChatType })
  type: ChatType;

  @ApiProperty({ example: 'Команда проекта', required: false })
  title?: string;

  @ApiProperty({ example: new Date().toISOString() })
  createdAt: string;

  @ApiProperty({ type: [Number], example: [1, 2, 3] })
  participantIds: number[];
}
