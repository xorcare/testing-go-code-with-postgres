CREATE TABLE users
(
    user_id    uuid,
    username   VARCHAR(50) UNIQUE       NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
