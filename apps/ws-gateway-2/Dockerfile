FROM node:20-slim

WORKDIR /app

COPY package*.json ./
COPY tsconfig*.json ./

COPY ./libs ./libs
COPY ./apps/ws-gateway-2 ./apps/ws-gateway-2

RUN npm install
RUN npm run build

CMD ["node", "dist/apps/ws-gateway-2/main.js"]