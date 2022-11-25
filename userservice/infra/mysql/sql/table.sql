CREATE DATABASE IF NOT EXISTS test;

USE test;

-- --------------------------
-- Table structure for users;
-- --------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键自增',
  `uid` int(11) UNSIGNED NOT NULL COMMENT 'uid',
  `username` varchar(32) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码(md5)',
  `sex` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别(0-男性 1-女性)',
  `register_time` timestamp(0) NULL DEFAULT NULL COMMENT '注册时间',
  `created_at` timestamp(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp(0) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uix_uid`(`uid`),
  UNIQUE INDEX `uix_username`(`username`)
) ENGINE=InnoDB CHARACTER SET = utf8mb4 COMMENT = '用户数据表';
