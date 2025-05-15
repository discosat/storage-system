INSERT INTO observation (id, observation_request_id, user_id, object_reference)
VALUES (1, 1, 1, 'object_reference_1');

INSERT INTO observation_metadata (id, measurement_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_cource)
VALUES (1, 1, 123456789, 1080, 1920, 3, 1746527580, 8, 0, 'Camera1', ST_SetSRID(ST_MakePoint(-37.968750, 74.683250), 4326), 1746527580, 1746527580, 50.0, 100.0, 180.0);