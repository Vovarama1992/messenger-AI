import { ApiProperty } from '@nestjs/swagger';

export class CreateMessageDto {
  @ApiProperty({ example: 1 })
  chatId: number;

  @ApiProperty({ example: 42 })
  senderId: number;

  @ApiProperty({ example: 'Привет!' })
  text: string;
}
