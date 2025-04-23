import { ApiProperty } from '@nestjs/swagger';

export class UserDto {
  @ApiProperty({ example: 1 })
  id: number;

  @ApiProperty({ example: 'user@example.com' })
  email: string;

  @ApiProperty({ example: 'Иван' })
  name: string;

  @ApiProperty({ example: new Date().toISOString() })
  createdAt: string;
}
