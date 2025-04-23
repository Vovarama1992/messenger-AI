FROM node:20-slim

WORKDIR /app

COPY package*.json ./
COPY tsconfig*.json ./

COPY ./libs ./libs
COPY ./apps/ws-gateway-ai-advice ./apps/ws-gateway-ai-advice

RUN npm install
RUN npm run build

CMD ["node", "dist/apps/ws-gateway-ai-advice/main.js"]