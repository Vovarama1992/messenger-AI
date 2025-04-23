import { Controller, Get, Request } from '@nestjs/common';
import { UsersService } from './users.service';
import { Request as Req } from 'express';
import { ApiResponse, ApiTags } from '@nestjs/swagger';

@Controller('users')
@ApiTags('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Get('me')
  @ApiResponse({
    status: 200,
    description: 'Successfully fetched the current user information.',
    example: {
      id: 1,
      email: 'user@example.com',
    },
  })
  async getMe(@Request() req: Req) {
    return this.usersService.authenticate(req);
  }
}
