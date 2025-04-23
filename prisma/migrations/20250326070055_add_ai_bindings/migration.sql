-- CreateTable
CREATE TABLE "AiChatBinding" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "chatId" INTEGER NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT "AiChatBinding_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "AiChatBinding_userId_chatId_key" ON "AiChatBinding"("userId", "chatId");

-- AddForeignKey
ALTER TABLE "AiChatBinding" ADD CONSTRAINT "AiChatBinding_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "AiChatBinding" ADD CONSTRAINT "AiChatBinding_chatId_fkey" FOREIGN KEY ("chatId") REFERENCES "Chat"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
