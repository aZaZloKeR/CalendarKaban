create table calendar.event
(
    id          bigserial not null primary key,
    date        date not null,
    time_start  time not null,
    time_end    time not null,
    name        varchar   not null,
    description varchar,
    user_id     int       not null references calendar."user"
)