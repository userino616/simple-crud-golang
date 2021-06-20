CREATE DATABASE crud_sample_db;

USE crud_sample_db;


CREATE TABLE user
(
    id      INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)
);

CREATE TABLE posts
(
    id      INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    title   VARCHAR(255),
    body    VARCHAR(1000),
    FOREIGN KEY (user_id) REFERENCES user (id)
);
