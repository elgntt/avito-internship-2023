CREATE TABLE segments (
    id              SERIAL PRIMARY KEY,
    slug            VARCHAR(255) UNIQUE
);

CREATE TABLE users (
    id              SERIAL PRIMARY KEY
);

CREATE TABLE users_segments (
    id              SERIAL PRIMARY KEY,
    user_id         INT,
    segment_id      INT,
    expiration_time TIMESTAMP WITH TIME ZONE,
    CONSTRAINT user_segment_unique UNIQUE (user_id, segment_id)
);