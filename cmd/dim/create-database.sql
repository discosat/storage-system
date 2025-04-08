CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    institution VARCHAR(255) NOT NULL
);

CREATE TABLE mission (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE,
    bucket VARCHAR(255) UNIQUE
);

CREATE TABLE request (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) UNIQUE,
     user_id INT NOT NULL,
     mission_id INT NOT NULL,
    --     created_at TIMESTAMP NOT NULL,
     CONSTRAINT fk_mission FOREIGN KEY(mission_id) REFERENCES mission(id),
     CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);

CREATE TABLE measurement_request(
    id SERIAL PRIMARY KEY,
    request_id INT NOT NULL,
    type VARCHAR(255),
    CONSTRAINT fk_request FOREIGN KEY(request_id) REFERENCES request(id)
);

CREATE TABLE observation (
    id SERIAL PRIMARY KEY,
    request_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_request FOREIGN KEY(request_id) REFERENCES request(id),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);

CREATE TABLE measurement (
    id SERIAL PRIMARY KEY,
    object_reference VARCHAR(255) UNIQUE,
    observation_id INT NOT NULL,
    measurement_request_id INT NOT NULL,
    CONSTRAINT fk_observation FOREIGN KEY(observation_id) REFERENCES observation(id),
    CONSTRAINT fk_measurement_request FOREIGN KEY(measurement_request_id) REFERENCES measurement_request(id)
);

CREATE TABLE measurement_metadata (
    id SERIAL PRIMARY KEY,
    measurement_id INT NOT NULL,
    location GEOGRAPHY(Point, 4326),
    CONSTRAINT fk_measurement FOREIGN KEY(measurement_id) REFERENCES measurement(id)
);
