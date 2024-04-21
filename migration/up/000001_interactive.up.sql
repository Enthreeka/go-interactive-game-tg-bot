set timezone = 'Europe/Moscow';

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
            CREATE TYPE role AS ENUM ('user', 'admin','superAdmin');
        END IF;
    END $$;

create table if not exists "user"(
    id           bigint unique,
    tg_username  text not null ,
    created_at   timestamp not null,
    phone        varchar(20) null,
    channel_from varchar(150) null,
    user_role         role default 'user' not null,
    blocked_bot bool default false,
    primary key (id)
);



-- todo create index for name
create table if not exists contest(
    id int generated always as identity,
    name varchar(200),
    file_id varchar(100),
    deadline timestamp,
    is_completed boolean default false,
    primary key (id)
);


create table if not exists questions(
    id int generated always as identity,
    contest_id int,
    created_by_user bigint,
    created_at timestamp,
    updated_at timestamp,
    question_name varchar(500),
    file_id varchar(100),
    deadline timestamp,
    is_send boolean default false,
    primary key (id),
    foreign key (contest_id)
        references contest (id) on delete cascade
);

create table if not exists answers(
    id int generated always as identity,
    answer varchar(100),
    cost_of_response int,
    primary key (id)
);

create table if not exists questions_answers(
    questions_id int,
    answers_id int,
    contest int,
    primary key (questions_id,answers_id),
    foreign key (questions_id)
        references questions (id) on delete cascade,
    foreign key (answers_id)
        references answers (id) on delete cascade,
    foreign key (contest)
        references contest (id) on delete cascade
);

create table if not exists user_results(
    id int generated always as identity,
    user_id bigint,
    contest_id int,
    total_points int,
    primary key (id),
    foreign key (user_id)
        references "user" (id) on delete cascade,
    foreign key (contest_id)
        references contest (id) on delete cascade
);

create table if not exists history_points(
    user_id bigint,
    questions_id int,
    awarded_point int,
    primary key (user_id,questions_id),
    foreign key (user_id)
        references "user" (id) on delete cascade,
    foreign key (questions_id)
        references questions (id) on delete cascade
);


update user_results set total_points = 0 where contest_id = 6;