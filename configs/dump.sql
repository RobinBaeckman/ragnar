UPDATE mysql.user SET Host='%' WHERE Host='localhost' AND User='rolf';
GRANT ALL PRIVILEGES ON *.* TO 'rolf'@'%';
FLUSH PRIVILEGES;
CREATE DATABASE IF NOT EXISTS rolf_db;
CREATE DATABASE IF NOT EXISTS rolf_db_test;
USE rolf_db;
SET SQL_MODE = 'STRICT_ALL_TABLES';
CREATE TABLE IF NOT EXISTS `users` (
	`id` CHAR(36) UNIQUE NOT NULL,
	`email` VARCHAR(64) UNIQUE NOT NULL,
	`password` VARCHAR(64) NOT NULL,
	`first_name` VARCHAR(64) NOT NULL,
	`last_name` VARCHAR(64) NOT NULL,
	`role` VARCHAR(64) NOT NULL,
	`token` VARCHAR(200),
	`created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
);
INSERT INTO users(id, email, password, first_name, last_name, role)
VALUES ('903c8d85-a877-4c40-8fa5-aea27f644921', 'rolle@mail.com', '$2a$10$drmi5DEMKThAMmloPtCo1Oc/q9uOBkVPztNlTisLxDQF36eNXtCHC', 'rolle', 'baeckman', 'admin');

USE rolf_db_test;
SET SQL_MODE = 'STRICT_ALL_TABLES';
CREATE TABLE IF NOT EXISTS `users` (
	`id` CHAR(36) UNIQUE NOT NULL,
	`email` VARCHAR(64) UNIQUE NOT NULL,
	`password` VARCHAR(64) NOT NULL,
	`first_name` VARCHAR(64) NOT NULL,
	`last_name` VARCHAR(64) NOT NULL,
	`role` VARCHAR(64) NOT NULL,
	`token` VARCHAR(200),
	`created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
	`created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
);
INSERT INTO users(id, email, password, first_name, last_name, role)
VALUES ('903c8d85-a877-4c40-8fa5-aea27f644921', 'rolle@mail.com', '$2a$10$drmi5DEMKThAMmloPtCo1Oc/q9uOBkVPztNlTisLxDQF36eNXtCHC', 'rolle', 'baeckman', 'admin');
