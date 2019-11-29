/*
Navicat MySQL Data Transfer

Source Server         : 118.25.123.163
Source Server Version : 50723
Source Host           : 118.25.123.163:3306
Source Database       : tuuser

Target Server Type    : MYSQL
Target Server Version : 50723
File Encoding         : 65001

Date: 2019-11-27 15:01:15
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for f_users
-- ----------------------------
DROP TABLE IF EXISTS `f_users`;
CREATE TABLE `f_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `f_union_id` varchar(64) NOT NULL DEFAULT '',
  `name` varchar(64) NOT NULL DEFAULT '',
  `mobile` varchar(16) NOT NULL DEFAULT '',
  `email` varchar(128) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `expiration` int(11) NOT NULL DEFAULT '0',
  `valid` tinyint(1) NOT NULL DEFAULT '1',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `union_id_idx` (`f_union_id`),
  UNIQUE KEY `mobile_idx` (`mobile`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=324 DEFAULT CHARSET=utf8;
