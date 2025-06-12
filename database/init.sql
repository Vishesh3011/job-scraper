SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS linkedin_job_scraper_geo_locations;

CREATE TABLE IF NOT EXISTS `job_scraper_users` (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `job_scraper_linkedin_geo_locations` (
    `id` BINARY(16)  NOT NULL,
    `location` VARCHAR(50) NOT NULL,
    `geo_id` VARCHAR(50) NOT NULL,
    `state` VARCHAR(50) NOT NULL,
    `country` VARCHAR(50) NOT NULL,
    CONSTRAINT `job_scraper_geo_locations_pkey` PRIMARY KEY (`id`),
    CONSTRAINT `job_scraper_geo_locations_geo_id_key` UNIQUE (`geo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;