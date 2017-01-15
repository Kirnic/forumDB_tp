DROP DATABASE IF EXISTS `forumDB`;
CREATE DATABASE `forumDB` DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
USE `forumDB`;




SET @PREVIOUS_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `user`;
DROP TABLE IF EXISTS `thread`;
DROP TABLE IF EXISTS `subscription`;
DROP TABLE IF EXISTS `post`;
DROP TABLE IF EXISTS `forum`;
DROP TABLE IF EXISTS `follow`;


CREATE TABLE `follow` (
  `follower` varchar(150) NOT NULL,
  `following` varchar(150) NOT NULL,
  PRIMARY KEY (`follower`,`following`),
  KEY `idx_following_follower` (`following`,`follower`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `forum` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(150) NOT NULL,
  `short_name` varchar(150) NOT NULL,
  `user` varchar(150) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_short_name` (`short_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=289 DEFAULT CHARSET=utf8;


CREATE TABLE `post` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `message` text NOT NULL,
  `parent` int(11) DEFAULT NULL,
  `likes` int(11) NOT NULL DEFAULT '0',
  `dislikes` int(11) NOT NULL DEFAULT '0',
  `points` int(11) NOT NULL DEFAULT '0',
  `isApproved` tinyint(4) NOT NULL DEFAULT '0',
  `isDeleted` tinyint(4) NOT NULL DEFAULT '0',
  `isEdited` tinyint(4) NOT NULL DEFAULT '0',
  `isHighlighted` tinyint(4) NOT NULL DEFAULT '0',
  `isSpam` tinyint(4) NOT NULL DEFAULT '0',
  `forum` varchar(150) NOT NULL,
  `thread` int(11) NOT NULL,
  `user` varchar(150) NOT NULL,
  `first_path` int(11) NOT NULL DEFAULT '0',
  `last_path` varchar(150) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_forum_date` (`forum`,`date`) USING BTREE,
  KEY `idx_user_date` (`user`,`date`) USING BTREE,
  KEY `idx_thread_date` (`thread`,`date`) USING BTREE,
  KEY `idx_thread_first_path_last_path` (`thread`,`first_path`,`last_path`) USING BTREE,
  KEY `idx_forum_user` (`forum`,`user`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1000454 DEFAULT CHARSET=utf8;


CREATE TABLE `subscription` (
  `user` varchar(150) NOT NULL,
  `thread` int(11) NOT NULL,
  PRIMARY KEY (`user`,`thread`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `thread` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(150) NOT NULL,
  `date` datetime NOT NULL,
  `slug` varchar(150) NOT NULL,
  `message` text NOT NULL,
  `likes` int(11) NOT NULL DEFAULT '0',
  `dislikes` int(11) NOT NULL DEFAULT '0',
  `points` int(11) NOT NULL DEFAULT '0',
  `posts` int(11) NOT NULL DEFAULT '0',
  `isClosed` tinyint(4) NOT NULL DEFAULT '0',
  `isDeleted` tinyint(4) NOT NULL DEFAULT '0',
  `forum` varchar(150) NOT NULL,
  `user` varchar(150) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_forum_date` (`forum`,`date`) USING BTREE,
  KEY `idx_user_date` (`user`,`date`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10372 DEFAULT CHARSET=utf8;


CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(150) NOT NULL,
  `username` varchar(150) DEFAULT NULL,
  `name` varchar(150) DEFAULT NULL,
  `about` text,
  `isAnonymous` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`) USING BTREE,
  UNIQUE KEY `idx_name` (`name`,`email`) USING BTREE,
  KEY `idx_id_name` (`id`,`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=100287 DEFAULT CHARSET=utf8;




SET FOREIGN_KEY_CHECKS = @PREVIOUS_FOREIGN_KEY_CHECKS;


