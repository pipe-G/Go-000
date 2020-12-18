CREATE TABLE `avatar_objects` (
  `id` bigint(20) NOT NULL,
  `type` bigint(20) NOT NULL,
  `object_id` bigint(20) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `objects` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `uuid` varchar(255) DEFAULT NULL,
  `mime_type` varchar(255) DEFAULT NULL,
  `deleted` tinyint(1) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;


CREATE TABLE `user_action_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `type` varchar(32) DEFAULT NULL,
  `referrer_url` varchar(256) DEFAULT NULL,
  `user_agent` varchar(256) DEFAULT NULL,
  `client_ip` varchar(64) DEFAULT NULL,
  `memo` varchar(128) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `type_index` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8;


CREATE TABLE `users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password` varchar(64) NOT NULL,
  `email` varchar(128) DEFAULT NULL,
  `phone_number` varchar(32) DEFAULT NULL,
  `gender` varchar(16) DEFAULT NULL,
  `avatar` varchar(128) DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `type` int(11) DEFAULT '0',
  `rights` int(11) DEFAULT '0',
  `meta` varchar(256) DEFAULT '{}',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `last_login_at` datetime DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username_unique` (`username`),
  KEY `created_at_index` (`created_at`),
  KEY `email_index` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1000004 DEFAULT CHARSET=utf8;


CREATE TABLE `verify_codes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `verify_type` int(11) DEFAULT NULL,
  `device_type` int(11) DEFAULT NULL,
  `to` varchar(64) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `verify_code` varchar(32) DEFAULT NULL,
  `expire_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8;
