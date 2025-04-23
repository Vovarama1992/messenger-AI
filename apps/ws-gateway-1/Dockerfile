FROM node:20-slim

WORKDIR /app

COPY package*.json ./
COPY tsconfig*.json ./

COPY ./libs ./libs
COPY ./apps/ws-gateway-1 ./apps/ws-gateway-1

RUN npm install
RUN npm run build

CMD ["node", "dist/apps/ws-gateway-1/main.js"]