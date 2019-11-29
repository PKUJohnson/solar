
DROP TABLE IF EXISTS `solar_users`;
CREATE TABLE `solar_users` (
  `app_id` varchar(128) DEFAULT NULL,
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `session_key` varchar(255) DEFAULT NULL,
  `open_id` varchar(255) DEFAULT NULL,
  `union_id` varchar(255) DEFAULT NULL,
  `user_name` varchar(255) DEFAULT NULL,
  `user_display_name` varchar(255) DEFAULT NULL,
  `gender` varchar(255) DEFAULT NULL,
  `mobile` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `province` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `language` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `signature` varchar(1024) NOT NULL DEFAULT '',
  `expiration` int(11) NOT NULL DEFAULT '0',
  `f_union_id` varchar(64) NOT NULL DEFAULT '',
  `avatar` varchar(512) NOT NULL DEFAULT '',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`user_id`),
  KEY `app_idx` (`app_id`),
  KEY `open_id_idx` (`open_id`),
  KEY `f_union_id_idx` (`f_union_id`)
) ENGINE=InnoDB AUTO_INCREMENT=18810 DEFAULT CHARSET=utf8;
