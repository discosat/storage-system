---------------------------------IMPORTANT-------------------------------------
-- Only ever add to this file don't delete, or you will have to adjust tests --
-------------------------------------------------------------------------------
INSERT INTO mission(name, bucket) VALUES ('Disco2', 'testbucket');

-- 1
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 1', 1, 1);
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (1, 1, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (2, 1, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (3, 1, 'image_series');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (4, 1, 'other');
--
-- 2
INSERT INTO flight_plan (name, mission_id, user_id, locked) VALUES ('flight plan 2', 1, 1, true);
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (5, 2, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (6, 2, 'other');
INSERT INTO observation(observation_request_id, user_id, object_reference, bucket_name) VALUES (6, 1, 'testDir/testFile.txt', 'testbucket');
INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES
    (1,12345678,1080,1920,2,123456789,6,24,'W', st_setsrid(st_point(10.4058633, 73), 4326),123456789,123456789,420,17000,2);
-- 3
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan 3', 1, 1);
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (7, 3, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (8, 3, 'other');
-- 4
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan update test', 1, 1);
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (40, 4, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (41, 4, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (42, 4, 'image_series');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (43, 4, 'other');
-- 5
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('flight plan update delete test', 1, 1);
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (432, 5, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (433, 5, 'image');

-- DAM TEST DATA
-- 6
INSERT INTO flight_plan (name, mission_id, user_id) VALUES ('Dam test', 1, 1);
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (1000, 6, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (1001, 6, 'image');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (1002, 6, 'image_series');
INSERT INTO observation_request(id, flight_plan_id, type) VALUES (1003, 6, 'image');

INSERT INTO observation(observation_request_id, user_id, object_reference, bucket_name) VALUES (1000, 1, 'Ooo.jpg', 'testbucket');
INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES
    (2,12345678,1080,1920,2,123456789,6,24,'N', st_setsrid(st_point(10.4058633, 73), 4326),123456789,123456789,420,17000,2);

INSERT INTO observation(observation_request_id, user_id, object_reference, bucket_name) VALUES (6, 1, 'Mary.jpg', 'testbucket');
INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES
    (3,12345678,1080,1920,2,123456789,6,24,'W', st_setsrid(st_point(10.4058633, 73), 4326),123456789,123456789,420,17000,2);

INSERT INTO observation(observation_request_id, user_id, object_reference, bucket_name) VALUES (6, 1, 'goblinMode.jpg', 'testbucket');
INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES
    (4,12345678,1080,1920,2,123456789,6,24,'I', st_setsrid(st_point(10.4058633, 73), 4326),123456789,123456789,420,17000,2);

INSERT INTO observation(observation_request_id, user_id, object_reference, bucket_name) VALUES (6, 1, 'forgor.jpg', 'testbucket');
INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES
    (5,12345678,1080,1920,2,123456789,6,24,'N', st_setsrid(st_point(10.4058633, 73), 4326),123456789,123456789,420,17000,2);

