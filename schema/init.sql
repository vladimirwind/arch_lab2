create table if not exists users(
    user_id int auto_increment,
    user_name varchar(255) not null,
    user_surname varchar(255) not null,
    user_login varchar(10) not null,
    user_password varchar(10) not null,
    primary key(user_id)
);

create table if not exists conferences(
    conference_id int auto_increment,
    conference_name varchar(255) not null,
    primary key(conference_id)
);

create table if not exists reports(
    report_id int auto_increment,
    report_name varchar(255) not null,
    user_id int,
    conference_id int,
    primary key(report_id),
    constraint fk_report_author foreign key (user_id) references users(user_id),
    constraint fk_conference foreign key (conference_id) references conferences(conference_id)
);

insert into users (user_name, user_surname, user_login, user_password) VALUES ('admin','admin','log1','pass1');
insert into reports (report_name, user_id) VALUES ("arch_lab_kursach", 1);
insert into conferences (conference_name) VALUES ("arch_lab_conference");