set time zone 'UTC';
create extension pgcrypto;

CREATE TABLE products (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);

INSERT INTO products (id, name, description) VALUES ('2310d6be-0e80-11ed-861d-0242ac120002', 'Red Panda', '8 Eastern Himalayas Drive');