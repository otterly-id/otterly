import { Elysia } from "elysia";
import { websocketRoute } from "./router/websocket.route";
import cors from "@elysiajs/cors";
import swagger from "@elysiajs/swagger";

const app = new Elysia()
    .use(cors({
      origin: "*",
      methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
      credentials: true,
    }))
    .use(swagger({
      path: "/docs",
      documentation: {
        info: {
          title: "Otterly MCP Host",
          version: "0.1.0",
        },
        tags: [
          {
            name: "Agent",
            description: "Agent related endpoints",
          },
        ],
      },
    }))
    .get("/", () => "Otterly MCP Host v0.1.0")
    .use(websocketRoute)
    .listen(3000);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
