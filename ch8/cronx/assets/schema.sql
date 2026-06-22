CREATE DATABASE IF NOT EXISTS cronx
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

USE cronx;

CREATE TABLE IF NOT EXISTS `job` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL DEFAULT '' COMMENT '任务名称',
  `namespace` varchar(45) NOT NULL DEFAULT '' COMMENT 'Kubernetes Namespace',
  `info` TEXT NOT NULL COMMENT 'Kubernetes Job 相关信息',
  `status` varchar(45) NOT NULL DEFAULT '' COMMENT '任务状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name_namespace` (`name`, `namespace`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

-- test data
INSERT IGNORE INTO `job` (`id`, `name`, `namespace`, `info`, `status`) VALUES (1, 'demo-job-1', 'cronx', '{"image":"alpine","command":["sleep"],"args":["60"]}', 'Normal');
INSERT IGNORE INTO `job` (`id`, `name`, `namespace`, `info`, `status`) VALUES (2, 'demo-job-2', 'cronx', '{"image":"busybox","command":["echo"],"args":["Hello Cronx!"]}', 'Normal');
