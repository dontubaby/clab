-- Удаление таблиц, если они существуют
DROP TABLE IF EXISTS leagues, users, resources, user_resources, neutrals, buildings, heroes, abilities, 
hero_ability, units, enemies, areas, areas_neutrals, areas_buildings, areas_heroes, 
areas_units, areas_enemies, actions, world_state;

-- Таблица leagues (лиги)
CREATE TABLE leagues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    authority INTEGER DEFAULT 0
);

-- Таблица users (пользователи)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    subscription BOOLEAN DEFAULT FALSE,
    league_id BIGINT REFERENCES leagues(id) ON DELETE SET NULL,
    balance DECIMAL(19, 2) DEFAULT 0,
    level INTEGER DEFAULT 1
);

-- Таблица ресурсов
CREATE TABLE resources (
  ID SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  value DECIMAL(19, 2) DEFAULT 0
);

-- Таблица связей user_resources (ресурсы пользователей)
CREATE TABLE user_resources (
    resource_id BIGINT NOT NULL REFERENCES resources(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
  	PRIMARY KEY  (resource_id,user_id)    
);

-- Создание таблицы neutrals (нейтральные объекты на арене)
CREATE TABLE neutrals (
    id BIGSERIAL PRIMARY KEY,                    
    name VARCHAR(255) NOT NULL,                  
    product VARCHAR(255) NOT NULL,                
    productivity_coefficient INTEGER NOT NULL,
    capacity DECIMAL(19, 2) NOT NULL,
    threshold_level1 DECIMAL(19, 2) NOT NULL,     
    threshold_level2 DECIMAL(19, 2) NOT NULL,
    size INTEGER NOT NULL,                            
    coordinates JSONB NOT NULL                    
);

-- Таблица buildings (здания)
CREATE TABLE buildings (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    product VARCHAR(255) NOT NULL,
    characteristics JSONB NOT NULL,
    level INTEGER DEFAULT 1,
    upgrade_price DECIMAL(19, 2) DEFAULT 0,
    size INTEGER NOT NULL,
    coordinates JSONB NOT NULL 
);

-- Таблица abilities (способности героев)
CREATE TABLE abilities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    characteristics JSONB NOT NULL,
    level INTEGER DEFAULT 1
 );

-- Таблица героев
CREATE TABLE heroes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    characteristics JSONB NOT NULL,
    experience DECIMAL(19, 2) DEFAULT 0,
    experience_to_up DECIMAL(19, 2) DEFAULT 0,
    level INTEGER DEFAULT 1,
     coordinates JSONB NOT NULL 
);

-- Таблица hero_ability (связь героев и способностей)
CREATE TABLE hero_ability (
  hero_id BIGINT NOT NULL REFERENCES heroes(id) ON DELETE CASCADE,
  ability_id BIGINT NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
  PRIMARY KEY (hero_id,ability_id)
);

-- Таблица units (юниты)
CREATE TABLE units (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    characteristics JSONB NOT NULL,
    experience DECIMAL(19, 2) DEFAULT 0,
    experience_to_up DECIMAL(19, 2) DEFAULT 0,
    level INTEGER DEFAULT 1,
    coordinates JSONB NOT NULL 
);

-- Таблица enemies (враги)
CREATE TABLE enemies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    characteristics JSONB NOT NULL,
    level INTEGER DEFAULT 1,
    coordinates JSONB NOT NULL 
);

-- Таблица areas (арены)
CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    width INTEGER NOT NULL CHECK(width > 0),
    height INTEGER NOT NULL CHECK(height > 0),
    cell_type_id INTEGER NOT NULL CHECK(cell_type_id BETWEEN 1 AND 4)   
);

-- Таблица связей арены и нейтральных объектов
CREATE TABLE areas_neutrals (
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    neutral_id BIGINT NOT NULL REFERENCES neutrals(id) ON DELETE CASCADE,
    PRIMARY KEY (area_id, neutral_id)
);

-- Таблица связей арены и зданий
CREATE TABLE areas_buildings (
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    building_id BIGINT NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    PRIMARY KEY (area_id, building_id)
);

-- Таблица связей арены и героев
CREATE TABLE areas_heroes (
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    hero_id BIGINT NOT NULL REFERENCES heroes(id) ON DELETE CASCADE,
    PRIMARY KEY (area_id, hero_id)
);

-- Таблица связей арены и юнитов
CREATE TABLE areas_units (
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    unit_id BIGINT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    PRIMARY KEY (area_id, unit_id)
);

-- Таблица связей арены и врагов
CREATE TABLE areas_enemies (
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    enemy_id BIGINT NOT NULL REFERENCES enemies(id) ON DELETE CASCADE,
    PRIMARY KEY (area_id, enemy_id)
);

-- Таблица actions (действия)
CREATE TABLE actions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    object_source_id BIGINT,
    object_dest_id BIGINT,
    action_type SMALLINT NOT NULL CHECK(action_type BETWEEN 1 AND 4),
    start_time TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    duration INTERVAL DEFAULT '00:00:00',
    status VARCHAR(255) DEFAULT ''  --Предусмотрено 3 статуса 1 - DONE, 2- NOT DONE , 3- PROCESS
);

-- Таблица world_state (состояние мира)
CREATE TABLE world_state (
    session_id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    area_id BIGINT NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    action_id BIGINT NOT NULL REFERENCES actions(id) ON DELETE CASCADE,
    time_stamp TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_user_resources_user_id ON user_resources(user_id);
CREATE INDEX idx_actions_user_id ON actions(user_id);

-- Наполнение таблицы leagues (лиги)
INSERT INTO leagues (name, authority) VALUES
('Bronze League', 100),
('Silver League', 200),
('Gold League', 300),
('Platinum League', 400);

-- Наполнение таблицы users (пользователи)
INSERT INTO users (login, password, email, subscription, league_id, balance, level) VALUES
('user1', 'password1', 'user1@example.com', TRUE, 1, 1000.00, 5),
('user2', 'password2', 'user2@example.com', FALSE, 2, 500.00, 3),
('user3', 'password3', 'user3@example.com', TRUE, 3, 750.00, 7),
('user4', 'password4', 'user4@example.com', FALSE, 4, 300.00, 2);

-- Наполнение таблицы resources (ресурсы)
INSERT INTO resources (name, value) VALUES
('Gold', 0),
('Wood', 0),
('Stone', 0),
('Food', 0);

-- Наполнение таблицы user_resources (ресурсы пользователей)
INSERT INTO user_resources (resource_id, user_id) VALUES
(1, 1), (2, 1), (3, 1), (4, 1),
(1, 2), (2, 2),
(1, 3), (3, 3),
(1, 4), (4, 4);

-- Наполнение таблицы neutrals (нейтральные объекты)
INSERT INTO neutrals (name, product, productivity_coefficient, capacity, threshold_level1, threshold_level2, size, coordinates) VALUES
('Forest', 'Wood', 10, 1000, 100, 500, 5, '[{"q": 10, "r": 20}]'),
('Gold Mine', 'Gold', 5, 500, 50, 250, 3, '[{"q": 30, "r": 40}]'),
('Stone Quarry', 'Stone', 8, 800, 80, 400, 4, '[{"q": 50, "r": 60}]'),
('Long Lake', 'Fish', 15, 1500, 200, 1000, 3, '[{"q": 30, "r": 40}, {"q": 31, "r": 40}, {"q": 32, "r": 40}]');

-- Наполнение таблицы buildings (здания)
INSERT INTO buildings (name, product, characteristics, level, upgrade_price, size, coordinates) VALUES
('Town Hall', 'Gold', '{"hp": 1000, "defense": 50}', 1, 500.00, 10, '[{"q": 70, "r": 80}]'),
('Barracks', 'Units', '{"hp": 800, "defense": 30}', 1, 300.00, 8, '[{"q": 90, "r": 100}]'),
('Farm', 'Food', '{"hp": 600, "defense": 20}', 1, 200.00, 6, '[{"q": 110, "r": 120}]'),
('Large Castle', 'Gold', '{"hp": 2000, "defense": 100}', 1, 1000.00, 4, '[{"q": 10, "r": 20}, {"q": 11, "r": 20}, {"q": 10, "r": 21}, {"q": 11, "r": 21}]');

-- Наполнение таблицы abilities (способности)
INSERT INTO abilities (name, characteristics, level) VALUES
('Fireball', '{"damage": 100, "cooldown": 5}', 1),
('Heal', '{"healing": 50, "cooldown": 10}', 1),
('Shield', '{"defense": 20, "duration": 15}', 1);

-- Наполнение таблицы heroes (герои)
INSERT INTO heroes (name, characteristics, experience, experience_to_up, level, coordinates) VALUES
('Hero1', '{"hp": 500, "attack": 50}', 0, 100, 1, '[{"q": 130, "r": 140}]'),
('Hero2', '{"hp": 600, "attack": 60}', 50, 150, 1, '[{"q": 150, "r": 160}]'),
('Hero3', '{"hp": 700, "attack": 70}', 100, 200, 1, '[{"q": 170, "r": 180}]');

-- Наполнение таблицы hero_ability (связь героев и способностей)
INSERT INTO hero_ability (hero_id, ability_id) VALUES
(1, 1), (1, 2),
(2, 2), (2, 3),
(3, 1), (3, 3);

-- Наполнение таблицы units (юниты)
INSERT INTO units (name, characteristics, experience, experience_to_up, level, coordinates) VALUES
('Warrior', '{"hp": 200, "attack": 20}', 0, 50, 1, '[{"q": 190, "r": 200}]'),
('Archer', '{"hp": 150, "attack": 30}', 0, 50, 1, '[{"q": 210, "r": 220}]'),
('Mage', '{"hp": 100, "attack": 40}', 0, 50, 1, '[{"q": 230, "r": 240}]'),
('Big Warrior', '{"hp": 300, "attack": 40}', 0, 100, 1, '[{"q": 50, "r": 60}, {"q": 50, "r": 61}]');

-- Наполнение таблицы enemies (враги)
INSERT INTO enemies (name, characteristics, level, coordinates) VALUES
('Goblin', '{"hp": 100, "attack": 10}', 1, '[{"q": 250, "r": 260}]'),
('Orc', '{"hp": 200, "attack": 20}', 2, '[{"q": 270, "r": 280}]'),
('Dragon', '{"hp": 500, "attack": 50}', 5, '[{"q": 290, "r": 300}]');

-- Наполнение таблицы areas (арены)
INSERT INTO areas (user_id, width, height, cell_type_id) VALUES
(1, 100, 100, 1),
(2, 100, 100, 2),
(3, 100, 100, 3),
(4, 100, 100, 4);

-- Наполнение таблицы areas_neutrals (связь арен и нейтральных объектов)
INSERT INTO areas_neutrals (area_id, neutral_id) VALUES
(1, 1), (1, 2),
(2, 2), (2, 3),
(3, 1), (3, 3);

-- Наполнение таблицы areas_buildings (связь арен и зданий)
INSERT INTO areas_buildings (area_id, building_id) VALUES
(1, 1), (1, 2),
(2, 2), (2, 3),
(3, 1), (3, 3);

-- Наполнение таблицы areas_heroes (связь арен и героев)
INSERT INTO areas_heroes (area_id, hero_id) VALUES
(1, 1), (1, 2),
(2, 2), (2, 3),
(3, 1), (3, 3);

-- Наполнение таблицы areas_units (связь арен и юнитов)
INSERT INTO areas_units (area_id, unit_id) VALUES
(1, 1), (1, 2),
(2, 2), (2, 3),
(3, 1), (3, 3);

-- Наполнение таблицы areas_enemies (связь арен и врагов)
INSERT INTO areas_enemies (area_id, enemy_id) VALUES
(1, 1), (1, 2),
(2, 2), (2, 3),
(3, 1), (3, 3);

-- Наполнение таблицы actions (действия)
INSERT INTO actions (user_id, area_id, object_source_id, object_dest_id, action_type, start_time, duration, status) VALUES
(1, 1, 1, 2, 1, CURRENT_TIMESTAMP, '00:05:00', 'PROCESS'),
(2, 2, 2, 3, 2, CURRENT_TIMESTAMP, '00:10:00', 'PROCESS'),
(3, 3, 3, 1, 3, CURRENT_TIMESTAMP, '00:15:00', 'PROCESS');

-- Наполнение таблицы world_state (состояние мира)
INSERT INTO world_state (user_id, area_id, action_id, time_stamp) VALUES
(1, 1, 1, CURRENT_TIMESTAMP),
(2, 2, 2, CURRENT_TIMESTAMP),
(3, 3, 3, CURRENT_TIMESTAMP);