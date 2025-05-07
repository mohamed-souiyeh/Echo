CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY, -- SERIAL is a common way for auto-incrementing integers in Postgres
    "username" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP -- TIMESTAMPTZ is preferred for timestamps
);

CREATE INDEX "username_idx" ON "users" ("username");
