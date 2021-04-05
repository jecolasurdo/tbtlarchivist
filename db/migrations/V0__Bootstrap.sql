CREATE TABLE `curated_clips` (
  `clip_id` int(11) NOT NULL AUTO_INCREMENT,
  `initial_date_curated` datetime NOT NULL,
  `last_date_curated` datetime NOT NULL,
  `curator_info` varchar(50) NOT NULL,
  `title` varchar(250) NOT NULL,
  `description` longtext NOT NULL,
  `media_uri` varchar(2048) NOT NULL,
  `media_type` varchar(3) NOT NULL,
  `priority` int(11) NOT NULL,
  PRIMARY KEY (`clip_id`),
  UNIQUE KEY `title_UNIQUE` (`title`),
  UNIQUE KEY `media_uri_UNIQUE` (`media_uri`) USING HASH
);

CREATE TABLE `curated_episodes` (
  `episode_id` int(11) NOT NULL AUTO_INCREMENT,
  `initial_date_curated` datetime NOT NULL,
  `last_date_curated` datetime NOT NULL,
  `curator_info` varchar(50) NOT NULL,
  `date_aired` datetime NOT NULL,
  `title` varchar(250),
  `description` longtext NOT NULL,
  `media_uri` varchar(2048) NOT NULL,
  `media_type` varchar(3) NOT NULL,
  `priority` int(11) NOT NULL,
  PRIMARY KEY (`episode_id`),
  UNIQUE KEY `date_aired_title_UNIQUE` (`date_aired`, `title`),
  UNIQUE KEY `media_uri_UNIQUE` (`media_uri`) USING HASH
);

CREATE TABLE `research_backlog` (
  `research_id` int(11) NOT NULL AUTO_INCREMENT,
  `episode_id` int(11) NOT NULL,
  `clip_id` int(11) NOT NULL,
  PRIMARY KEY (`research_id`),
  UNIQUE KEY `episode_id_clip_id_UNIQUE` (`episode_id`, `clip_id`)
);

CREATE TABLE `research_leases` (
  `lease_id` char(36) NOT NULL,
  `research_id` int(11) NOT NULL,
  `expiration` datetime NOT NULL,
  PRIMARY KEY (`lease_id`, `research_id`),
  UNIQUE KEY `research_id_UNIQUE` (`research_id`)
);

CREATE TABLE `research_complete` (
  `research_id` int(11) NOT NULL,
  `episode_id` int(11) NOT NULL,
  `clip_id` int(11) NOT NULL,
  `episode_duration_ns` bigint(20) NOT NULL,
  `clip_duration_ns` bigint(20) NOT NULL,
  `research_date` datetime NOT NULL,
  PRIMARY KEY (`research_id`),
  UNIQUE KEY `episode_id_clip_id_UNIQUE` (`episode_id`, `clip_id`)
);

CREATE TABLE `episode_clip_offsets` (
  `research_id` int(11) NOT NULL,
  `offset_ns` bigint(20) NOT NULL,
  PRIMARY KEY (`research_id`)
);

CREATE TABLE `episode_hashes` (
  `episode_id` int(11) NOT NULL,
  `hash` varchar(32) NOT NULL,
  PRIMARY KEY (`episode_id`)
);

CREATE INDEX episode_hash_idx ON episode_hashes(hash);

CREATE TABLE `clip_hashes` (
  `clip_id` int(11) NOT NULL,
  `hash` varchar(32) NOT NULL,
  PRIMARY KEY (`clip_id`)
);

CREATE INDEX clip_hash_idx ON clip_hashes(hash);