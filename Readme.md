## To create Database
CREATE DATABASE tasksDB


## To create Table
CREATE TABLE tasks (
   id serial PRIMARY KEY,
   name VARCHAR(50) NOT NULL,
   content VARCHAR(255) NOT NULL
);


## To Populate Table
INSERT INTO tasks(name, content) VALUES('James', 'Dummy note'),('Karl', 'I must conquer Europe'),('Fer', 'Buy BTC')


## To Remove Table
DROP TABLE tasks
