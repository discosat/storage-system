INSERT INTO "user"(name, institution) VALUES ('Hans Pedersen', 'SDU');

INSERT INTO mission(name, bucket) VALUES ('Disco2', 'disco2data');

INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 1', 1, 1);
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 2', 1, 1);
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 3', 1, 1);

-- flightPlan observation requests list
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'Wide-Image');
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'Narrow-Image');
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'Thermal-Image');
INSERT INTO observation_request(flight_plan_id, type) VALUES (1, 'Other');
-- Expand list when needed
