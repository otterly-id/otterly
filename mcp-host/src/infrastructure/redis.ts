import { createClient, RedisClientType } from "redis";

export const redis: RedisClientType = createClient({
    url: process.env.REDIS_URL,
});

export async function initRedis() {
    await redis.connect();
    console.log("Redis connected");
}