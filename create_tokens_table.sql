create table tokens
(
    id         integer                 not null
        constraint tokens_pk
            primary key,
    user_id    integer,
    email      varchar(255)            not null,
    token      varchar(255)            not null,
    token_hash bytea                   not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    expiry     timestamp               not null
);

comment on table tokens is 'Table for storing user''s tokens';

comment on column tokens.id is 'ID of the token';

comment on column tokens.user_id is 'User''s ID';

comment on column tokens.email is 'User''s email';

comment on column tokens.token is 'Token';

comment on column tokens.token_hash is 'Hash of the token';

comment on column tokens.created_at is 'Created At';

comment on column tokens.updated_at is 'Updated At';

comment on column tokens.expiry is 'Token''s expiry date';

alter table tokens
    owner to postgres;

alter table public.tokens
    add constraint tokens_user_id_fk
        foreign key (user_id) references public.users
            on update cascade on delete cascade;

comment on constraint tokens_user_id_fk on public.tokens is 'Foreign Key to Users table';

