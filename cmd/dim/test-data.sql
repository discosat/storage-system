INSERT INTO "user"(name, institution) VALUES ('Hans Pedersen', 'SDU');

INSERT INTO mission(name, bucket) VALUES ('Disco2', 'disco2data');

INSERT INTO request (name, mission_id, user_id) VALUES ('observation 1', 1, 1);

-- Measurement request list
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Wide-Image');
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Narrow-Image');
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Thermal-Image');
INSERT INTO measurement_request(request_id, type) VALUES (1, 'Other');
-- Expand list when needed

-- SELECT * FROM request r
--     INNER JOIN public.measurement_request mr on r.id = mr.request_id WHERE r.id = 1;
--
-- SELECT * FROM measurement_request where request_id = 1 AND type = 'Wide-Image';

-- SELECT * FROM request r
--     INNER JOIN public.measurement_request mr on r.id = mr.request_id WHERE r.id = 1;
--
-- SELECT * FROM measurement_request where request_id = 1 AND type = 'Wide-Image';
-- SELECT * FROM request r
--     INNER JOIN public.measurement_request mr on r.id = mr.request_id WHERE r.id = 1;
--
-- SELECT * FROM measurement_request where request_id = 1 AND type = 'Wide-Image';
