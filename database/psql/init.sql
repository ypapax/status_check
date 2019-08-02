CREATE DATABASE status_check;
\connect status_check;
CREATE TABLE "service"
(
    id      bigserial primary key,
    address varchar(100) NOT NULL UNIQUE,
    created timestamp    NOT NULL
);

CREATE TABLE "status"
(
    id          bigserial primary key,
    available   boolean   NOT NULL,
    created     timestamp NOT NULL,
    response_ms integer   NOT NULL,
    service_id  integer REFERENCES service (id)
);

CREATE INDEX status_created_idx ON status USING btree(created, available, response_ms);