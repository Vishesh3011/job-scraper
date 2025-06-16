SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS job_scraper_geo_locations;

INSERT INTO job_scraper_linkedin_geo_locations (id, location, geo_id, state, country) VALUES
    (UUID_TO_BIN(UUID()), 'Sydney', '104769905', 'New South Wales', 'AU'),
    (UUID_TO_BIN(UUID()), 'New South Wales', '103313686', 'New South Wales', 'AU'),
    (UUID_TO_BIN(UUID()), 'Canberra', '106089960', 'Australian Capital Territory', 'AU'),
    (UUID_TO_BIN(UUID()), 'Melbourne', '100992797', 'Victoria', 'AU'),
    (UUID_TO_BIN(UUID()), 'Perth', '103392068', 'Western Australia', 'AU'),
    (UUID_TO_BIN(UUID()), 'Adelaide', '107042567', 'South Australia', 'AU'),
    (UUID_TO_BIN(UUID()), 'Brisbane', '104468365', 'Queensland', 'AU'),
;

SET FOREIGN_KEY_CHECKS = 1;