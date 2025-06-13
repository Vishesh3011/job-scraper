SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS job_scraper_geo_locations;

INSERT INTO job_scraper_linkedin_geo_locations (id, location, geo_id, state, country) VALUES
    (UUID_TO_BIN(UUID()), 'Sydney', '104769905', 'New South Wales', 'AU'),
    (UUID_TO_BIN(UUID()), 'Melbourne', '104769906', 'Victoria', 'AU'),
    (UUID_TO_BIN(UUID()), 'Brisbane', '104769907', 'Queensland', 'AU'),
    (UUID_TO_BIN(UUID()), 'Perth', '104769908', 'Western Australia', 'AU'),
    (UUID_TO_BIN(UUID()), 'Adelaide', '104769909', 'South Australia', 'AU'),
    (UUID_TO_BIN(UUID()), 'Canberra', '104769910', 'Australian Capital Territory', 'AU'),
    (UUID_TO_BIN(UUID()), 'Hobart', '104769911', 'Tasmania', 'AU'),
    (UUID_TO_BIN(UUID()), 'Darwin', '104769912', 'Northern Territory', 'AU'),
    (UUID_TO_BIN(UUID()), 'Auckland', '104769913', 'Auckland', 'NZ'),
    (UUID_TO_BIN(UUID()), 'Wellington', '104769914', 'Wellington', 'NZ'),
    (UUID_TO_BIN(UUID()), 'Christchurch', '104769915', 'Canterbury', 'NZ'),
    (UUID_TO_BIN(UUID()), 'Hamilton', '104769916', 'Waikato', 'NZ'),
    (UUID_TO_BIN(UUID()), 'Dunedin', '104769917', 'Otago', 'NZ')
;

SET FOREIGN_KEY_CHECKS = 1;