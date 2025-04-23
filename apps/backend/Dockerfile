FROM node:20

WORKDIR /app

COPY package*.json ./
COPY tsconfig*.json ./

COPY ./libs ./libs
COPY ./prisma ./prisma               
COPY ./apps/backend ./apps/backend

RUN npm install
RUN npm run build

# 👇 Добавляем миграции перед стартом
CMD npx prisma migrate deploy && node dist/apps/backend/main.js