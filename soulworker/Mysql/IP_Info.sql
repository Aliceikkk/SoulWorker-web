/*
 Navicat Premium Data Transfer

 Source Server         : IP数据库
 Source Server Type    : MySQL
 Source Server Version : 50740
 Source Host           : 47.242.190.82:3306
 Source Schema         : IP_Info

 Target Server Type    : MySQL
 Target Server Version : 50740
 File Encoding         : 65001

 Date: 18/05/2024 15:03:30
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for Ban_ip
-- ----------------------------
DROP TABLE IF EXISTS `Ban_ip`;
CREATE TABLE `Ban_ip`  (
  `IP` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `TIME` int(30) NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of Ban_ip
-- ----------------------------

-- ----------------------------
-- Table structure for Info
-- ----------------------------
DROP TABLE IF EXISTS `Info`;
CREATE TABLE `Info`  (
  `IP` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `NUM` int(30) NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of Info
-- ----------------------------
INSERT INTO `Info` VALUES ('115.201.226.236', 1);
INSERT INTO `Info` VALUES ('42.244.62.134', 1);

SET FOREIGN_KEY_CHECKS = 1;
