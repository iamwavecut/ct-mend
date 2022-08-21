create table if not exists clients
(
    id integer not null
    constraint clients_pk
    primary key autoincrement,
    name text not null,
    code_scan_interval int
);

create table if not exists projects
(
    id integer not null
    constraint projects_pk
    primary key autoincrement,
    client_id integer,
    name text
);