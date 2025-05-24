-- -- Create Database
-- CREATE OR REPLACE FUNCTION set_updated_at()
--     RETURNS TRIGGER AS
-- $$
-- BEGIN
--     NEW.updated_at = NOW();
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TABLE mission
-- (
--     id         SERIAL PRIMARY KEY,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     name       VARCHAR(50) UNIQUE,
--     bucket     VARCHAR(255) UNIQUE
-- );
--
-- CREATE TRIGGER updated_at_trigger
--     BEFORE UPDATE
--     ON mission
--     FOR EACH ROW
-- EXECUTE PROCEDURE set_updated_at();
--
--
-- CREATE TABLE flight_plan
-- (
--     id         SERIAL PRIMARY KEY,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     name       VARCHAR(255) UNIQUE,
--     user_id    INT                                    NOT NULL,
--     mission_id INT                                    NOT NULL,
--     CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES mission (id)
-- );
--
-- CREATE TRIGGER updated_at_trigger
--     BEFORE UPDATE
--     ON flight_plan
--     FOR EACH ROW
-- EXECUTE PROCEDURE set_updated_at();
--
-- CREATE TYPE observation_type AS ENUM ('image', 'image_series', 'number', 'other');
--
-- CREATE TABLE observation_request
-- (
--     id             SERIAL PRIMARY KEY,
--     created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     flight_plan_id INT                                    NOT NULL,
--     type           observation_type,
--     CONSTRAINT fk_flight_plan FOREIGN KEY (flight_plan_id) REFERENCES flight_plan (id)
-- );
--
-- CREATE TABLE observation_request_metadata
-- (
--     id                     SERIAL PRIMARY KEY,
--     created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     updated_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     observation_request_id INT                                    NOT NULL,
--     metadata               VARCHAR(255),
--     CONSTRAINT fk_observation_request FOREIGN KEY (observation_request_id) REFERENCES observation_request (id)
-- );
--
-- CREATE TRIGGER updated_at_trigger
--     BEFORE UPDATE
--     ON observation_request_metadata
--     FOR EACH ROW
-- EXECUTE PROCEDURE set_updated_at();
--
--
-- CREATE TABLE observation
-- (
--     id                     SERIAL PRIMARY KEY,
--     created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     updated_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     observation_request_id INT                                    NOT NULL,
--     user_id                INT                                    NOT NULL,
--     object_reference       VARCHAR(255) UNIQUE                    NOT NULL,
--     bucket_name            VARCHAR(255)                           NOT NULL,
--     CONSTRAINT fk_request FOREIGN KEY (observation_request_id) REFERENCES observation_request (id)
-- );
--
--
-- -- START Bucket enforcing
-- CREATE OR REPLACE FUNCTION update_observation_buckets()
--     RETURNS TRIGGER AS
-- $$
-- BEGIN
--     UPDATE observation o
--     SET bucket_name = NEW.bucket
--     FROM observation_request orq
--              JOIN flight_plan fp ON orq.flight_plan_id = fp.id
--     WHERE o.observation_request_id = orq.id
--       AND fp.mission_id = NEW.id;
--
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER update_observation_buckets_trigger
--     AFTER UPDATE OF bucket
--     ON mission
--     FOR EACH ROW
--     WHEN (OLD.bucket IS DISTINCT FROM NEW.bucket)
-- EXECUTE FUNCTION update_observation_buckets();
-- -- END bucket enforcing
--
--
-- CREATE TRIGGER updated_at_trigger
--     BEFORE UPDATE
--     ON observation
--     FOR EACH ROW
-- EXECUTE PROCEDURE set_updated_at();
--
-- CREATE TABLE observation_metadata
-- (
--     id             SERIAL PRIMARY KEY,
--     created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
--     observation_id INT UNIQUE                             NOT NULL,
--     size           INT,
--     height         INT,
--     width          INT,
--     channels       INT,
--     timestamp      INT,
--     bits_pixels    INT,
--     image_offset   INT,
--     camera         VARCHAR(255),
--     location       GEOGRAPHY(Point, 4326),
--     gnss_date      INT,
--     gnss_time      INT,
--     gnss_speed     FLOAT,
--     gnss_altitude  FLOAT,
--     gnss_cource    FLOAT,
--     CONSTRAINT fk_observation FOREIGN KEY (observation_id) REFERENCES observation (id)
-- );
--
-- CREATE TRIGGER updated_at_trigger
--     BEFORE UPDATE
--     ON observation_metadata
--     FOR EACH ROW
-- EXECUTE PROCEDURE set_updated_at();
-- -- END create database

INSERT INTO mission(name, bucket) VALUES ('Disco2', 'disco2data');

INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 1', 1, 1);
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 2', 1, 1);
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 3', 1, 1);

-- flightPlan observation requests list
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'image');
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'image');
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'image_series');
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'other');
-- Expand list when needed


INSERT INTO observation_request(flight_plan_id, type) VALUES (2, 'image');
INSERT INTO observation_request(flight_plan_id, type) VALUES (2, 'other');
INSERT INTO observation(observation_request_id, user_id, object_reference, bucket_name) VALUES (6, 1, 'testDir/testFile.txt', 'testBucket');
INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES
    (1,12345678,1080,1920,2,123456789,6,24,'W', st_setsrid(st_point(10.4058633, 73), 4326),123456789,123456789,420,17000,2);

