CREATE TABLE IF NOT EXISTS taxi_trips (
    trip_id VARCHAR(255),
    trip_type VARCHAR(255),
    taxi_id VARCHAR(255),
    trip_start_timestamp TIMESTAMP WITH TIME ZONE,
    trip_end_timestamp TIMESTAMP WITH TIME ZONE,
    pickup_centroid_latitude DOUBLE PRECISION,
    pickup_centroid_longitude DOUBLE PRECISION,
    pickup_community_area INTEGER,
    dropoff_centroid_latitude DOUBLE PRECISION,
    dropoff_centroid_longitude DOUBLE PRECISION,
    dropoff_community_area INTEGER,
    pickup_zip_code VARCHAR(255),
    dropoff_zip_code VARCHAR(255),
    CONSTRAINT pk_trip PRIMARY KEY (trip_id, trip_type)
);

CREATE TABLE IF NOT EXISTS covid_daily (
    report_date DATE,
    update_date DATE,
    cases INTEGER,
    deaths INTEGER,
    hospitalizations INTEGER,
    CONSTRAINT pk_covid_daily PRIMARY KEY (report_date)
);

CREATE TABLE IF NOT EXISTS covid_weekly_zip (
    zip_code INTEGER,
    week_number INTEGER,
    week_start DATE,
    week_end DATE,
    update_date DATE,
    cases INTEGER,
    cases_rate REAL,
    tests INTEGER,
    tests_rate REAL,
    percent_tested_positive REAL,
    deaths INTEGER,
    deaths_rate REAL,
    CONSTRAINT pk_covid_weekly PRIMARY KEY (zip_code, week_start, week_end)
);

CREATE TABLE IF NOT EXISTS ccvi_ca (
    community_area INTEGER,
    community_area_name VARCHAR(255),
    ccvi_score REAL,
    ccvi_category VARCHAR(255),
    CONSTRAINT pk_ccvi_ca PRIMARY KEY (community_area)
);

CREATE TABLE IF NOT EXISTS ccvi_zip (
    zip_code INTEGER,
    ccvi_score REAL,
    ccvi_category VARCHAR(255),
    CONSTRAINT pk_ccvi_zip PRIMARY KEY (zip_code)
);

CREATE TABLE IF NOT EXISTS building_permits (
    id VARCHAR(255),
    permit_nbr INTEGER,
    permit_type VARCHAR(255),
    review_type VARCHAR(255),
    application_date DATE,
    issue_date DATE,
    total_fee REAL,
    zip_code INTEGER,
    community_area INTEGER,
    CONSTRAINT pk_building PRIMARY KEY (id, permit_nbr)
);

CREATE TABLE IF NOT EXISTS unemployment (
    community_area INTEGER,
    community_area_name VARCHAR(255),
    unemployment_rate REAL,
    poverty_rate REAL,
    per_capita_income INTEGER,
    CONSTRAINT pk_unemployment PRIMARY KEY (community_area)
);

CREATE TABLE IF NOT EXISTS community_area_names (
    community_area INTEGER,
    community_area_name VARCHAR(255),
    CONSTRAINT pk_community_area_names PRIMARY KEY (community_area)
);

CREATE TABLE IF NOT EXISTS zip2ca (
    zip_code INTEGER,
    community_area1 INTEGER,
    community_area2 INTEGER,
    community_area3 INTEGER,
    community_area4 INTEGER,
    community_area5 INTEGER,
    community_area6 INTEGER,
    community_area7 INTEGER,
    community_area8 INTEGER,
    community_area9 INTEGER,
    CONSTRAINT pk_zip2ca PRIMARY KEY (zip_code)
);
