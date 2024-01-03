
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY NOT NULL,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
	name VARCHAR(50) NOT NULL DEFAULT '',
	status SMALLINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX tasks_uuid_idx ON tasks (uuid);

-- +migrate Down
DROP TABLE IF EXISTS tasks;
