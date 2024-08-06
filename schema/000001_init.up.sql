CREATE TABLE customers
(
    id serial not null unique,
    name varchar(255) not null,
    surname varchar(255) not null,
    age int not null,
    balance int not null,

    phone varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE storage(
    storage int
);
INSERT INTO storage (storage) VALUES (0);

CREATE TABLE credits
(
    id serial not null unique,
    customer_id int not null,
    value int not null,
    percentage int not null,
    loan_period int not null,
    current_debt int not null
);