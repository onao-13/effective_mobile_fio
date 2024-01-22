create table if not exists humans
(
    id         bigserial not null
        constraint table_name_pk
            primary key,
    name       varchar(50)                                           not null,
    surname    varchar(50)                                           not null,
    patronymic varchar(50)                                           not null,
    age        integer                                               not null,
    gender     varchar(6)                                            not null
);

create table if not exists humans_nationality
(
    humanid     bigint      not null
        constraint humans_nationality_humans_id_fk
            references humans,
    countryid   varchar(10) not null,
    probability real        not null
);

create index human__surname__index
    on humans (surname);

alter table humans_nationality
    drop constraint humans_nationality_humans_id_fk;

alter table humans_nationality
    add constraint humans_nationality_humans_id_fk
        foreign key (humanid) references humans
            on delete cascade;
