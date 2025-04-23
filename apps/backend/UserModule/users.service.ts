import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from 'libs/PrismaModule/prisma.service';
import { Request } from 'express';

import { ConfigService } from '@nestjs/config';

@Injectable()
export class UsersService {
  private readonly baseurl: string;

  constructor(
    private prisma: PrismaService,
    private jwtService: JwtService,
    private configService: ConfigService,
  ) {
    this.baseurl =
      this.configService.get<string>('BASE_URL') ||
      'https://app.opticard.co/api';
  }

  async authenticate(req: Request) {
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      throw new UnauthorizedException(
        'Authorization header for route getMe missing',
      );
    }

    const token = authHeader.split(' ')[1];
    if (!token) {
      throw new UnauthorizedException('Token missing');
    }

    try {
      const decoded = this.jwtService.verify(token, {
        secret: process.env.JWT_SECRET,
      });
      const userId = decoded.sub;

      const user = await this.prisma.user.findUnique({
        where: { id: userId },
      });

      if (!user) {
        throw new UnauthorizedException('User not found');
      }

      const { password, ...rest } = user;

      console.log(password);

      return rest;
    } catch (error) {
      throw new UnauthorizedException(error.message);
    }
  }
}
