
-- +migrate Up
CREATE TABLE IF NOT EXISTS "notifications"
(
    "id"                    VARCHAR(255) NOT NULL DEFAULT '',
    "from"                  VARCHAR(255) NOT NULL DEFAULT '',
    "to"                    VARCHAR(255) NOT NULL DEFAULT '',
    "content"               VARCHAR(255),
    "status"                VARCHAR(255) NOT NULL DEFAULT '',
    "created_at"            TIMESTAMPTZ DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ DEFAULT NOW(),
);

-- +migrate Down
DROP TABLE "notifications";