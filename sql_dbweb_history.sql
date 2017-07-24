CREATE TABLE `dbweb_history` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `user` varchar(50) DEFAULT NULL,
  `db` varchar(50) NOT NULL DEFAULT 'all',
  `tag` varchar(20) NOT NULL,
  `changes` mediumtext,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cretime` (`create_time`),
  KEY `idx_user_ctime` (`user`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
