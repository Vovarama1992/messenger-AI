import { Test, TestingModule } from '@nestjs/testing';
import { UsersController } from './users.controller';
import { UsersService } from './users.service';
import { JwtService } from '@nestjs/jwt';
import { PrismaService } from 'src/PrismaModule/prisma.service';
import { UnauthorizedException } from '@nestjs/common';

describe('UsersController', () => {
  let controller: UsersController;
  let service: UsersService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [UsersController],
      providers: [
        UsersService,
        {
          provide: JwtService,
          useValue: { verify: jest.fn() },
        },
        {
          provide: PrismaService,
          useValue: { user: { findUnique: jest.fn() } },
        },
      ],
    }).compile();

    controller = module.get<UsersController>(UsersController);
    service = module.get<UsersService>(UsersService);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });

  it('should authenticate and get user data', async () => {
    const mockRequest = { headers: { authorization: 'Bearer valid_token' } };
    const mockUser = { id: 1, email: 'user@example.com', role: 'PRIVATE_PERSON' };

    jest.spyOn(service, 'authenticate').mockResolvedValue(mockUser); // mock authenticate method

    const result = await controller.getMe(mockRequest as any);

    expect(result).toEqual(mockUser);
    expect(service.authenticate).toHaveBeenCalledWith(mockRequest);
  });

  it('should throw UnauthorizedException when no token is provided', async () => {
    const mockRequest = { headers: {} };

    await expect(controller.getMe(mockRequest as any)).rejects.toThrow(UnauthorizedException);
  });
});