SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';


DROP DATABASE IF EXISTS `sspanel`;
CREATE DATABASE IF NOT EXISTS `sspanel`;

GRANT ALL PRIVILEGES ON sspanel.* TO 'sspanel'@'%';
DROP USER 'sspanel'@'%';

CREATE USER 'sspanel'@'%' IDENTIFIED BY 'sspanel';
GRANT ALL PRIVILEGES ON sspanel.* TO 'sspanel'@'%';

USE sspanel;


DROP TABLE IF EXISTS `ss_node`;
CREATE TABLE `ss_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `enableota` TINYINT(1) NOT NULL DEFAULT '1',
  `server` varchar(128) NOT NULL,
  `method` varchar(64) NOT NULL DEFAULT 'aes-256-cfb',
  `custom_method` tinyint(1) NOT NULL DEFAULT '0',
  `traffic_rate` tinyint(4) NOT NULL DEFAULT '1',
  `description` varchar(128) NOT NULL DEFAULT '',
  `status` TINYINT(1) NOT NULL DEFAULT '1',
  `traffic_limit` bigint(63) NOT NULL DEFAULT '0',
  `upload` bigint(63) NOT NULL DEFAULT '0',
  `download` bigint(63) NOT NULL DEFAULT '0',
  `total_upload` bigint(63) NOT NULL DEFAULT '0',
  `total_download` bigint(63) NOT NULL DEFAULT '0',
  `location`  varchar(128) NOT NULL DEFAULT '',
  `vps_server_id` varchar(128) NOT NULL DEFAULT '0',
  `vps_server_name` varchar(128) NOT NULL DEFAULT '',
  UNIQUE KEY `server` (`server`, `name`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(128) CHARACTER SET utf8mb4 NOT NULL,
  `email` varchar(32) NOT NULL,
  `manage_pass` varchar(64) NOT NULL,
  `passwd` varchar(64) NOT NULL,
  `traffic_rate` float NOT NULL DEFAULT '1',
  `upload` bigint(63) NOT NULL DEFAULT '0',
  `download` bigint(63) NOT NULL DEFAULT '0',
  `traffic_limit` bigint(63) NOT NULL DEFAULT '0',
  `total_upload` bigint(63) NOT NULL DEFAULT '0',
  `total_download` bigint(63) NOT NULL DEFAULT '0',
  `enable_ota` tinyint(4) NOT NULL DEFAULT '1',
  `last_check_in_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_reset_pass_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reg_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_admin` int(2) NOT NULL DEFAULT '0',
  `expire_time` timestamp NOT NULL,
  `is_email_verify` tinyint(4) NOT NULL DEFAULT '0',
  `reg_ip` varchar(128) NOT NULL DEFAULT '127.0.0.1',
  `description` varchar(256) NOT NULL,
  `status` TINYINT(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user` (`user_name`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `user` (`user_name`, `email`, `manage_pass`, `passwd`, `is_admin`, `description`, `expire_time`) VALUES ('admin', 'admin@gmail.com', 'f4c6b8435ce61e80e30bed9d0c28d5d9', '33629385d2e0b1d5b4b8566c50be1552', '1', 'admin', '2030-12-01 22:01:23');

DROP TABLE IF EXISTS `user_token`;
CREATE TABLE `user_token` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `token` varchar(256) NOT NULL,
  `user_id` int(11) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `expire_time` timestamp NOT NULL,
  `name` varchar(128) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `userkey` (`user_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



DROP TABLE IF EXISTS `api_server`;
CREATE TABLE `api_server` (
  `id` int(32) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `host` varchar(128) NOT NULL,
  `port` int(32) NOT NULL,
  `status` TINYINT(1)  NOT NULL DEFAULT '1',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `host` (`host`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `vps_server_account`;
CREATE TABLE `vps_server_account` (
  `id` int(32) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `operators` varchar(64) NOT NULL,
  `api_key`   varchar(128) NOT NULL,
  `credit_ceilings`   decimal(19,4)  NOT NULL,
  `lables`    varchar(128),
  `expire_time`  timestamp NOT NULL,
  `descryption` text NOT NULL,
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




DELIMITER |

CREATE EVENT IF NOT EXISTS `sspanel`.`expireToken`
ON SCHEDULE EVERY 1 HOUR
  DO
    BEGIN
      DELETE FROM `sspanel`.`user_token` WHERE `expire_time` < NOW();
    END

| DELIMITER ;


DELIMITER |

CREATE TRIGGER IF NOT EXISTS `sspanel`.`sumUserTraffic`  BEFORE UPDATE ON user
ON SCHEDULE EVERY 1 HOUR
  DO
    BEGIN
      DELETE FROM `sspanel`.`user_token` WHERE `expire_time` < NOW();
    END

| DELIMITER ;
