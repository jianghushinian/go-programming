DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id`            int(11) NOT NULL AUTO_INCREMENT,
  `uuid`          varchar(50)           DEFAULT '' COMMENT 'UUID',
  `name`          varchar(50)           DEFAULT '' COMMENT '用户名',
  `email`         varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
  `age`           tinyint(4) NOT NULL DEFAULT '0' COMMENT '年龄',
  `birthday`      datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '生日',
  `member_number` varchar(50) COMMENT '成员编号',
  `activated_at`  datetime COMMENT '激活时间',
  `created_at`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`    datetime,
  PRIMARY KEY (`id`),
  UNIQUE KEY `u_email` (`email`),
  INDEX           `idx_deleted_at`(`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `id`         int(11) NOT NULL AUTO_INCREMENT,
  `title`      varchar(100)      DEFAULT '' COMMENT '标题',
  `content`    text COMMENT '内容',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime          DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY          `idx_deleted_at`(`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章表';

DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id`         int(11) NOT NULL AUTO_INCREMENT,
  `content`    text COMMENT '评论',
  `post_id`    int(11) COMMENT '文章 ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime          DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY          `idx_deleted_at`(`deleted_at`),
  KEY          `fk_post_comments`(`post_id`),
  CONSTRAINT `fk_post_comments` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章评论表';

DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags` (
  `id`         int(11) NOT NULL AUTO_INCREMENT,
  `name`       varchar(100)      DEFAULT '' COMMENT '标签名',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime          DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY          `idx_deleted_at`(`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章标签表';

DROP TABLE IF EXISTS `posts_tags`;
CREATE TABLE `posts_tags` (
  `post_id` int(11) NOT NULL COMMENT '文章 ID',
  `tag_id`  int(11) NOT NULL COMMENT '标签 ID',
  PRIMARY KEY (`post_id`, `tag_id`),
  KEY       `fk_posts_tags_post` (`post_id`),
  CONSTRAINT `fk_posts_tags_post` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
  CONSTRAINT `fk_posts_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章标签关联表';
