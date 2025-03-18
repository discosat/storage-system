INSERT INTO "user"(name, institution) VALUES ('Hans Pedersen', 'SDU') RETURNING id;

INSERT INTO mission(name, bucket) VALUES ('Disco2', 'disco2data') RETURNING id;

INSERT INTO request (name, mission_id, user_id) VALUES ('observation 1', 1, 1) RETURNING id;

-- Measurement request list
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Wide-Image') RETURNING id;
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Narrow-Image') RETURNING id;
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Thermal-Image') RETURNING id;
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Other') RETURNING id;
-- Expand list when needed

