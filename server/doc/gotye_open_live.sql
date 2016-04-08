/*
Source Server         : 192.168.1.10
Source Server Version : 50627
Source Host           : 192.168.1.10:3306
Source Database       : gotye_open_live

Target Server Type    : MYSQL
Target Server Version : 50627
File Encoding         : 65001

Date: 2016-04-07 15:29:41
*/

-- ----------------------------
-- Current Database: `gotye_open_live`
-- ----------------------------

CREATE DATABASE `gotye_open_live`;

USE `gotye_open_live`;


-- ----------------------------
-- Table structure for tbl_app
-- ----------------------------
DROP TABLE IF EXISTS `tbl_app`;
CREATE TABLE `tbl_app` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(64) NOT NULL,
  `value` varchar(128) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tbl_follow_liverooms
-- ----------------------------
DROP TABLE IF EXISTS `tbl_follow_liverooms`;
CREATE TABLE `tbl_follow_liverooms` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `liveroom_id` bigint(20) NOT NULL,
  `follow_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_user_liveroom_id` (`user_id`,`liveroom_id`) USING BTREE,
  KEY `index_user_id` (`user_id`) USING BTREE,
  KEY `index_liveroom_id` (`liveroom_id`) USING BTREE,
  CONSTRAINT `foreign_liveroom_id` FOREIGN KEY (`liveroom_id`) REFERENCES `tbl_liverooms` (`liveroom_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tbl_liverooms
-- ----------------------------
DROP TABLE IF EXISTS `tbl_liverooms`;
CREATE TABLE `tbl_liverooms` (
  `user_id` bigint(20) NOT NULL,
  `liveroom_id` bigint(20) NOT NULL,
  `liveroom_name` varchar(256) NOT NULL,
  `liveroom_desc` varchar(1024) DEFAULT NULL,
  `liveroom_topic` varchar(256) DEFAULT NULL,
  `anchor_pwd` varchar(64) NOT NULL,
  `assist_pwd` varchar(64) NOT NULL,
  `user_pwd` varchar(64) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`liveroom_id`),
  UNIQUE KEY `index_user_id` (`user_id`) USING BTREE,
  CONSTRAINT `foreign_user_id` FOREIGN KEY (`user_id`) REFERENCES `tbl_users` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tbl_online_liverooms
-- ----------------------------
DROP TABLE IF EXISTS `tbl_online_liverooms`;
CREATE TABLE `tbl_online_liverooms` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `liveroom_id` bigint(20) NOT NULL,
  `pushing_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `player_num` int(11) NOT NULL DEFAULT '0' COMMENT '用户观看人数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_liveroom_id` (`liveroom_id`) USING BTREE,
  UNIQUE KEY `index_id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tbl_pictures
-- ----------------------------
DROP TABLE IF EXISTS `tbl_pictures`;
CREATE TABLE `tbl_pictures` (
  `pic_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '图片id',
  `pic` blob NOT NULL,
  PRIMARY KEY (`pic_id`),
  UNIQUE KEY `index_pic_id` (`pic_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tbl_users
-- ----------------------------
DROP TABLE IF EXISTS `tbl_users`;
CREATE TABLE `tbl_users` (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '唯一账号',
  `phone` varchar(20) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `email` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `nickname` varchar(64) NOT NULL,
  `pwd` varchar(32) NOT NULL,
  `headpic_id` bigint(20) NOT NULL DEFAULT '0',
  `sex` enum('male','female') NOT NULL DEFAULT 'male' COMMENT '性别,默认男',
  `address` varchar(128) NOT NULL DEFAULT '中国,上海',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `index_userid` (`user_id``) USING BTREE,
  UNIQUE KEY `index_account` (`account`) USING BTREE,
  UNIQUE KEY `index_phone` (`phone`) USING BTREE,
  UNIQUE KEY `index_email` (`email`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
