create table calendar."user"
(
    id       bigserial not null primary key,
    username varchar   not null,
    email    varchar   not null unique,
    password varchar   not null
)