import { Elysia, t } from "elysia";
import { ChatSession, GoogleGenerativeAI } from "@google/generative-ai";
import { config } from "../config";

const genAI = new GoogleGenerativeAI(config.GEMINI_API_KEY);

type WebsocketMessage = {
    query: {
        apiKey: string;
    };
    chat?: ChatSession;
};

type WebsocketBody = {
    message: string;
}

type WebsocketStore = {
    chat?: ChatSession;
}

export const websocketRoute = new Elysia({
  websocket: {
    idleTimeout: 300,
    maxPayloadLength: 32 * 1024 * 1024,
  },
})
    .decorate('store', {} as WebsocketStore)
    .ws("/chat", {
      query: t.Object({
        apiKey: t.String(),
      }),
      body: t.Object({
        message: t.String(),
      }),
      open(ws) {
        const clientApiKey = ws.data.query.apiKey;
        const expectedApiKey = config.API_KEY;
        console.log(
          `Client API Key: '${clientApiKey}' (Type: ${typeof clientApiKey}, Length: ${clientApiKey.length})`
        );
        console.log(
          `Expected API Key: '${expectedApiKey}' (Type: ${typeof expectedApiKey}, Length: ${expectedApiKey.length})`
        );
        console.log(clientApiKey === expectedApiKey);

        // Validate API key
        if (!expectedApiKey || clientApiKey !== expectedApiKey) {
            console.warn(`[Auth] Invalid API Key. Disconnecting client ${ws.id}`);
            ws.close(1008, "Invalid API Key");
            return;
        }

        console.log(`[Open] Client connected: ${ws.id}`);

        try {
            const model = genAI.getGenerativeModel({
                model: "gemini-2.5-flash",
            })
            ws.data.store.chat = model.startChat();

            ws.send(JSON.stringify({
                type: "connection",
                status: "connected",
                message: "Chat session initialized"
              }));
        } catch (error) {
            console.error(`[Error] Failed to create chat session: ${error}`);
            ws.close(1008, "Failed to create chat session");
        }
      },
      async message(ws, message: WebsocketBody) {
        console.log(`[Message] From ${ws.id}: ${message}`);

        if (!ws.data.store.chat) {
            console.error(`[Error] No chat session for client ${ws.id}`);
            ws.send(JSON.stringify({
              type: "error",
              message: "Chat session not initialized"
            }));
            return;
          }

          // Validate message format
          if (!message || typeof message.message !== 'string' || message.message.trim().length === 0) {
            ws.send(JSON.stringify({
              type: "error",
              message: "Invalid message format"
            }));
            return;
          }

          try {
            const chat = ws.data.store.chat;
            
            // Send typing indicator
            if (ws.raw.readyState === 1) {
              ws.send(JSON.stringify({
                type: "typing",
                status: "start"
              }));
            }
      
            // Get streamed response from Gemini
            const result = await chat.sendMessageStream(message.message.trim());
            
            // Stream response chunks
            let fullResponse = "";
            for await (const chunk of result.stream) {
              if (ws.raw.readyState !== 1) {
                console.log(`[Warning] Connection closed during streaming for ${ws.id}`);
                break;
              }
      
              const chunkText = chunk.text();
              if (chunkText) {
                fullResponse += chunkText;
                ws.send(JSON.stringify({
                  type: "chunk",
                  content: chunkText
                }));
              }
            }

            // Send final response
            ws.send(JSON.stringify({
              type: "response",
              content: fullResponse
            }));
          } catch (error) {
            console.error(`[Error] Failed to send message: ${error}`);
            ws.send(JSON.stringify({
              type: "error",
              message: "Failed to send message"
            }));
          }
      },
      close(ws, code, reason) {
        console.log(`[Close] Client disconnected: ${ws.id}`, { code, reason });
        if (ws.data.store.chat) {
            ws.data.store.chat = undefined;
        }
      },
    });