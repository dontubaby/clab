DROP IF EXISTS world_state, users, leagues, user_resources, areas, area_objects, buildings, units, enemies, actions;

CREATE TABLE world_state (
    session_id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    area_id BIGINT NOT NULL REFERENCES areas(id),
    action_id BIGINT NOT NULL REFERENCES actions(id),
    time_stamp TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    subscription BOOLEAN DEFAULT FALSE,
    league_id BIGINT REFERENCES leagues(id) DEFAULT 0,
     resources JSONB NOT NULL DEFAULT '[]',
    balance DECIMAL(19, 2) DEFAULT 0,
    level INTEGER DEFAULT 1
);

CREATE TABLE leagues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    authority INTEGER DEFAULT 0
);

CREATE TABLE user_resources (
    resource_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    value DECIMAL(19, 2) DEFAULT 0
);


CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    cell_type_id SMALLINT NOT NULL CHECK(cell_type_id BETWEEN 1 AND 4),
     objects JSONB NOT NULL DEFAULT '[]'
);

CREATE TABLE area_objects (
    object_id BIGINT PRIMARY KEY,
    area_id BIGINT NOT NULL REFERENCES areas(id),
    object_type_id SMALLINT NOT NULL CHECK(object_type_id BETWEEN 1 AND 3)
);

CREATE TABLE buildings (
    building_id SERIAL PRIMARY KEY,
    characteristics JSONB NOT NULL,
    level INTEGER DEFAULT 1
);

CREATE TABLE units (
    unit_id SERIAL PRIMARY KEY,
    characteristics JSONB NOT NULL,
    level INTEGER DEFAULT 1
);

CREATE TABLE enemies (
    enemy_id SERIAL PRIMARY KEY,
    characteristics JSONB NOT NULL,
    level INTEGER DEFAULT 1
);

CREATE TABLE actions (
    action_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    area_id BIGINT NOT NULL REFERENCES areas(id),
    object_source_id BIGINT REFERENCES area_objects(object_id),
    object_dest_id BIGINT REFERENCES area_objects(object_id),
    action_type SMALLINT NOT NULL CHECK(action_type BETWEEN 1 AND 4),
    start_time TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    duration INTERVAL DEFAULT '00:00:00',
    status BOOLEAN DEFAULT FALSE
);
