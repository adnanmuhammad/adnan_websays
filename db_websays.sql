/*
Navicat MySQL Data Transfer

Source Server         : Localhost
Source Server Version : 50505
Source Host           : localhost:3306
Source Database       : db_websays

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2023-06-05 22:25:21
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for tbl_products
-- ----------------------------
DROP TABLE IF EXISTS `tbl_products`;
CREATE TABLE `tbl_products` (
  `prod_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `prod_title` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`prod_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of tbl_products
-- ----------------------------
INSERT INTO `tbl_products` VALUES ('1', 'test product title 1 - Edited Again');
INSERT INTO `tbl_products` VALUES ('2', 'test product title 2');
INSERT INTO `tbl_products` VALUES ('3', 'test product title 3');
INSERT INTO `tbl_products` VALUES ('5', 'test product title 4');
