CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE mission
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    name       VARCHAR(50) UNIQUE,
    bucket     VARCHAR(255) UNIQUE
);

CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON mission
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();


CREATE TABLE flight_plan
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    name       VARCHAR(255) UNIQUE,
    user_id    INT                                    NOT NULL,
    mission_id INT                                    NOT NULL,
    CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES mission (id)
);

CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON flight_plan
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();


CREATE TABLE observation_request
(
    id             SERIAL PRIMARY KEY,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    flight_plan_id INT                                    NOT NULL,
    type           VARCHAR(255),
    CONSTRAINT fk_flight_plan FOREIGN KEY (flight_plan_id) REFERENCES flight_plan (id)
);

CREATE TABLE observation_request_metadata
(
    id                     SERIAL PRIMARY KEY,
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    observation_request_id INT                                    NOT NULL,
    metadata               VARCHAR(255),
    CONSTRAINT fk_observation_request FOREIGN KEY (observation_request_id) REFERENCES observation_request (id)
);

CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON observation_request_metadata
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();


CREATE TABLE observation
(
    id                     SERIAL PRIMARY KEY,
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    observation_request_id INT                                    NOT NULL,
    user_id                INT                                    NOT NULL,
    object_reference       VARCHAR(255) UNIQUE                    NOT NULL,
    bucket_name            VARCHAR(255)                           NOT NULL,
    CONSTRAINT fk_request FOREIGN KEY (observation_request_id) REFERENCES observation_request (id)
);


-- START Bucket enforcing
CREATE OR REPLACE FUNCTION set_observation_bucket()
    RETURNS TRIGGER AS
$$
BEGIN
    SELECT m.bucket
    INTO NEW.bucket
    FROM mission m
             JOIN flight_plan fp ON fp.mission_id = m.id
             JOIN observation_request orq ON orq.flight_plan_id = fp.id
    WHERE orq.id = NEW.observation_request_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_bucket_trigger
    BEFORE INSERT
    ON observation
    FOR EACH ROW
EXECUTE PROCEDURE set_observation_bucket();

CREATE OR REPLACE FUNCTION prevent_bucket_update()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.bucket IS DISTINCT FROM OLD.bucket THEN
        RAISE EXCEPTION 'Bucket field is immutable and cannot be updated';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_bucket_update_trigger
    BEFORE UPDATE
    ON observation
    FOR EACH ROW
EXECUTE FUNCTION prevent_bucket_update();

-- END bucket enforcing


CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON observation
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

CREATE TABLE observation_metadata
(
    id             SERIAL PRIMARY KEY,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    observation_id INT UNIQUE                             NOT NULL,
    size           INT,
    height         INT,
    width          INT,
    channels       INT,
    timestamp      INT,
    bits_pixels    INT,
    image_offset   INT,
    camera         VARCHAR(255),
    location       GEOGRAPHY(Point, 4326),
    gnss_date      INT,
    gnss_time      INT,
    gnss_speed     FLOAT,
    gnss_altitude  FLOAT,
    gnss_cource    FLOAT,
    CONSTRAINT fk_measurement FOREIGN KEY (observation_id) REFERENCES observation (id)
);

CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON observation_metadata
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
