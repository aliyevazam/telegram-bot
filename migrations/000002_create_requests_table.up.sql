CREATE TABLE IF NOT EXISTS "requests"(
    "tg_id" BIGINT REFERENCES users("tg_id"),
    "request" TEXT, 
    "request_time" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);