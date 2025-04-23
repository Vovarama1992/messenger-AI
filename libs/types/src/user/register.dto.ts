import { ApiProperty } from '@nestjs/swagger';

export class RegisterDto {
  @ApiProperty({ example: 'user@example.com' })
  email: string;

  @ApiProperty({ example: 'qwerty123' })
  password: string;

  @ApiProperty({ example: 'Иван' })
  name: string;
}
