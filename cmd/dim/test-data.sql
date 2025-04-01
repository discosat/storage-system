INSERT INTO "user"(name, institution) VALUES ('Hans Pedersen', 'SDU');

INSERT INTO mission(name, bucket) VALUES ('Disco2', 'disco2data');

INSERT INTO request (name, mission_id, user_id) VALUES ('observation 1', 1, 1);
INSERT INTO request (name, mission_id, user_id) VALUES ('observation 2', 1, 1);
INSERT INTO request (name, mission_id, user_id) VALUES ('observation 3', 1, 1);

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

SELECT r.name FROM measurement
    INNER JOIN public.measurement_request mr on mr.id = measurement.measurement_request_id
    INNER JOIN public.request r on r.id = mr.request_id;

SELECT r.name FROM measurement
    INNER JOIN observation o on o.id = measurement.observation_id
    INNER JOIN public.request r on o.request_id = r.id;

SELECT * FROM observation where request_id = 1;

SELECT * FROM request
                  FULL JOIN measurement_request mr on request.id = mr.request_id
                  FULL JOIN measurement m on mr.id = m.measurement_request_id WHERE request.id = 1 AND m.id IS NOT NULL;

SELECT r.name, m.object_reference FROM observation
                                           FULL JOIN measurement m on observation.id = m.observation_id
                                           FULL JOIN public.request r on r.id = observation.request_id WHERE request_id = 1;

SELECT * FROM measurement WHERE id = 6;

SELECT r.id, r.name, r.user_id, r.mission_id FROM request r
    LEFT JOIN public.observation o on r.id = o.request_id WHERE o.id IS NULL;