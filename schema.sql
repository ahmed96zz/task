create table profile (
    id SERIAL Primary key,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(25) NOT NULL UNIQUE,
    Address text);

 create table users (
    id SERIAL Primary key,
    email VARCHAR(255) ,
    phone_number VARCHAR(12) NOT NULL UNIQUE,
    otp VARCHAR(10) NOT NULL DEFAULT '',
    otp_expiration_time TIMESTAMP);

