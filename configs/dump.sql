UPDATE mysql.user SET Host='%' WHERE Host='localhost' AND User='ruser';
FLUSH PRIVILEGES;
CREATE DATABASE IF NOT EXISTS ragnar_db;
CREATE TABLE IF NOT EXISTS `users` (
	`id` CHAR(36) UNIQUE NOT NULL,
	`email` VARCHAR(64) UNIQUE NOT NULL,
	`password` VARCHAR(64) NOT NULL,
	`first_name` VARCHAR(64) NOT NULL,
	`last_name` VARCHAR(64) NOT NULL,
	`role` VARCHAR(64) NOT NULL,
	PRIMARY KEY (`id`)
);
