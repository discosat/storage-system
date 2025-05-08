CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE "user"
(
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    name        VARCHAR(255)                           NOT NULL,
    institution VARCHAR(255)                           NOT NULL
);

CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();


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
    CONSTRAINT fk_mission FOREIGN KEY (mission_id) REFERENCES mission (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id)
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
    object_reference       VARCHAR(255) UNIQUE,
    CONSTRAINT fk_request FOREIGN KEY (observation_request_id) REFERENCES observation_request (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id)
);


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
    measurement_id INT                                    NOT NULL,
    size           INT,
    height         INT,
    width           INT,
    channels       INT,
    timestamp       INT,
    bits_pixels     INT,
    image_offset    INT,
    camera          VARCHAR(255),
    location       GEOGRAPHY(Point, 4326),
    gnss_date       INT,
    gnss_time       INT,
    gnss_speed      FLOAT,
    gnss_altitude   FLOAT,
    gnss_cource     FLOAT,
    CONSTRAINT fk_measurement FOREIGN KEY (measurement_id) REFERENCES observation (id)
);

CREATE TRIGGER updated_at_trigger
    BEFORE UPDATE
    ON observation_metadata
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
