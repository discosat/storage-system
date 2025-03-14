CREATE TABLE mission (
    id SERIAL PRIMARY KEY,
    missionName VARCHAR(50) UNIQUE,
    bucket VARCHAR(250) UNIQUE
);

CREATE TABLE IF NOT EXISTS measurements (
    id SERIAL PRIMARY KEY,
    mission_id INT NOT NULL,
    ref VARCHAR(250) NOT NULL
);

