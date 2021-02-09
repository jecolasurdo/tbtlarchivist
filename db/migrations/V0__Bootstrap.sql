CREATE TABLE `curated_clips` (
  `clip_id` int(11) NOT NULL AUTO_INCREMENT,
  `date_curated` datetime NOT NULL,
  `curator_info` varchar(50) NOT NULL,
  `title` varchar(250) NOT NULL,
  `description` longtext NOT NULL,
  `media_uri` varchar(2048) NOT NULL,
  `media_type` varchar(3) NOT NULL,
  PRIMARY KEY (`clip_id`),
  UNIQUE KEY `title_UNIQUE` (`title`),
  UNIQUE KEY `media_uri_UNIQUE` (`media_uri`) USING HASH
);

CREATE TABLE `curated_episodes` (
  `episode_id` int(11) NOT NULL AUTO_INCREMENT,
  `date_curated` datetime NOT NULL,
  `curator_info` varchar(50) NOT NULL,
  `date_aired` datetime NOT NULL,
  `title` varchar(250),
  `description` longtext NOT NULL,
  `media_uri` varchar(2048) NOT NULL,
  `media_type` varchar(3) NOT NULL,
  PRIMARY KEY (`episode_id`),
  UNIQUE KEY `date_aired_title_UNIQUE` (`date_aired`, `title`),
  UNIQUE KEY `media_uri_UNIQUE` (`media_uri`) USING HASH
);

CREATE TABLE `episode_leases` (
  `episode_id` int(11) NOT NULL,
  `expiration` datetime NOT NULL,
  PRIMARY KEY (`episode_id`)
);

CREATE TABLE `episode_clip_research` (
  `episode_clip_id` int(11) NOT NULL AUTO_INCREMENT,
  `episode_id` int(11) NOT NULL,
  `clip_id` int(11) NOT NULL,
  `episode_duration_ns` bigint(20) NOT NULL,
  `clip_duration_ns` bigint(20) NOT NULL,
  `research_date` datetime NOT NULL,
  PRIMARY KEY (`episode_clip_id`),
  UNIQUE KEY `episode_id_clip_id_UNIQUE` (`episode_id`, `clip_id`)
);

CREATE TABLE `episode_clip_offsets` (
  `episode_clip_id` int(11) NOT NULL,
  `offset_ns` bigint(20) NOT NULL,
  PRIMARY KEY (`episode_clip_id`)
);