DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `company`;

create table IF not exists `company`
(
    `id`               INT(20) AUTO_INCREMENT,
    `name`             VARCHAR(20) NOT NULL,
    `created_at`       Datetime DEFAULT current_timestamp,
    `updated_at`       Datetime DEFAULT current_timestamp on update current_timestamp,
    PRIMARY KEY (`id`)
    ) DEFAULT CHARSET=utf8;

create table IF not exists `user`
(
    `id`               INT(20) AUTO_INCREMENT,
    `cid`              INT(20) NOT NULL,
    `name`             VARCHAR(20) NOT NULL,
    `created_at`       Datetime DEFAULT current_timestamp,
    `updated_at`       Datetime DEFAULT current_timestamp on update current_timestamp,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`cid`) REFERENCES company(`id`)
    ) DEFAULT CHARSET=utf8;
