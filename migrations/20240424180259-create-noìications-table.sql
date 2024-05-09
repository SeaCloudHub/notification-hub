
-- +migrate Up
CREATE TABLE IF NOT EXISTS "notifications"
(
    "id"                    VARCHAR(255) NOT NULL DEFAULT '',
    "from_user"                  VARCHAR(255) NOT NULL DEFAULT '',
    "to_user"                    VARCHAR(255) NOT NULL DEFAULT '',
    "content"               VARCHAR(255),
    "status"                VARCHAR(255) NOT NULL DEFAULT '',
    "created_at"            TIMESTAMPTZ DEFAULT NOW(),
    "updated_at"            TIMESTAMPTZ DEFAULT NOW(),
    "viewed_at"             TIMESTAMPTZ
);

-- +migrate Down
DROP TABLE "notifications";