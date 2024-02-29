create table users
(
    id         integer                 not null
        constraint users_pk
            primary key,
    email      varchar(255)            not null,
    first_name varchar(255)            not null,
    last_name  varchar(255)            not null,
    password   varchar(60)             not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

comment on table users is 'Table of users';

comment on column users.id is 'ID of user';

comment on column users.email is 'User''s email address';

comment on column users.first_name is 'User''s first name';

comment on column users.last_name is 'User''s last name';

comment on column users.password is 'Password hash';

comment on column users.created_at is 'Created At';

comment on column users.updated_at is 'Updated At';

alter table users
    owner to postgres;