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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
