CREATE DATABASE `gowiki` DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
CREATE TABLE `gowiki`.`articles` (
  `title` CHAR(255) NOT NULL,
  `body` VARCHAR(10000) DEFAULT NULL,
  PRIMARY KEY(`title`)
);
GRANT ALL ON `gowiki`.`articles` TO 'root'@'%' IDENTIFIED by 'root';
