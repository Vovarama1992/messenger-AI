import { ApiProperty } from '@nestjs/swagger';

export class MessageDto {
  @ApiProperty({ example: 1 })
  id: number;

  @ApiProperty({ example: 1 })
  chatId: number;

  @ApiProperty({ example: 42 })
  senderId: number;

  @ApiProperty({ example: 'Привет!' })
  text: string;

  @ApiProperty({ example: new Date().toISOString() })
  createdAt: string;
}
