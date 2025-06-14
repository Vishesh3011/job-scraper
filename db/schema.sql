CREATE TABLE job_scraper_users (
                                   id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
                                   name VARCHAR(100) NOT NULL,
                                   email VARCHAR(100) UNIQUE NOT NULL,
                                   location VARCHAR(100) NOT NULL,
                                   keywords TEXT NOT NULL,
                                   cookie TEXT,
                                   csrf_token VARCHAR(100),
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE job_scraper_linkedin_geo_locations (
                                                    id BINARY (16) NOT NULL,
                                                    location VARCHAR(50) NOT NULL,
                                                    geo_id VARCHAR(50) NOT NULL,
                                                    state VARCHAR(50) NOT NULL,
                                                    country VARCHAR(50) NOT NULL,
                                                    CONSTRAINT job_scraper_geo_locations_pkey PRIMARY KEY (id),
                                                    CONSTRAINT job_scraper_geo_locations_geo_id_key UNIQUE (geo_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;