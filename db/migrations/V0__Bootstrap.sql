CREATE TABLE `clip_info` (
  `clip_id` int(11) NOT NULL AUTO_INCREMENT,
  `date_curated` datetime NOT NULL,
  `curator_info` varchar(50) NOT NULL,
  `title` varchar(250) NOT NULL,
  `description` longtext NOT NULL,
  `media_uri` varchar(2048) NOT NULL,
  `media_type` varchar(20) NOT NULL,
  PRIMARY KEY (`clip_id`),
  UNIQUE KEY `title_UNIQUE` (`title`),
  UNIQUE KEY `media_uri_UNIQUE` (`media_uri`) USING HASH
);

CREATE TABLE `episode_info` (
  `episode_id` int(11) NOT NULL AUTO_INCREMENT,
  `date_curated` datetime NOT NULL,
  `curator_info` varchar(50) NOT NULL,
  `date_aired` datetime NOT NULL,
  `duration` int(11),
  `title` varchar(250),
  `description` longtext NOT NULL,
  `media_uri` varchar(2048) NOT NULL,
  `media_type` varchar(20) NOT NULL,
  PRIMARY KEY (`episode_id`),
  UNIQUE KEY `date_aired_title_UNIQUE` (`date_aired`, `title`),
  UNIQUE KEY `media_uri_UNIQUE` (`media_uri`) USING HASH
);

CREATE TABLE `episode_leases` (
  `episode_id` int(11) NOT NULL,
  `expiration` datetime NOT NULL,
  PRIMARY KEY (`episode_id`)
);