CREATE TABLE "users" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT, -- SQLite auto-increments INTEGER PRIMARY KEY
    "username" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "username_idx" ON "users" ("username");