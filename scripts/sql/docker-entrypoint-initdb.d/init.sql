-- -------------------------------- The script used when storeMode is 'db' --------------------------------

CREATE database if NOT EXISTS `seata` default character set utf8mb4 collate utf8mb4_unicode_ci;
USE `seata`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
-- the table to store GlobalSession data
CREATE TABLE IF NOT EXISTS `global_table`
(
    `addressing` varchar(128) NOT NULL,
    `xid` varchar(128) NOT NULL,
    `transaction_id` bigint DEFAULT NULL,
    `transaction_name` varchar(128) DEFAULT NULL,
    `timeout` int DEFAULT NULL,
    `begin_time` bigint DEFAULT NULL,
    `status` tinyint NOT NULL,
    `active` bit(1) NOT NULL,
    `gmt_create` datetime DEFAULT NULL,
    `gmt_modified` datetime DEFAULT NULL,
    PRIMARY KEY (`xid`),
    KEY `idx_gmt_modified_status` (`gmt_modified`,`status`),
    KEY `idx_transaction_id` (`transaction_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- the table to store BranchSession data
CREATE TABLE IF NOT EXISTS `branch_table`
(
    `addressing` varchar(128) NOT NULL,
    `xid` varchar(128) NOT NULL,
    `branch_id` bigint NOT NULL,
    `transaction_id` bigint DEFAULT NULL,
    `resource_id` varchar(256) DEFAULT NULL,
    `lock_key` VARCHAR(1000),
    `branch_type` varchar(8) DEFAULT NULL,
    `status` tinyint DEFAULT NULL,
    `application_data` varchar(2000) DEFAULT NULL,
    `async_commit` tinyint NOT NULL DEFAULT 0,
    `gmt_create` datetime(6) DEFAULT NULL,
    `gmt_modified` datetime(6) DEFAULT NULL,
    PRIMARY KEY (`branch_id`),
    KEY `idx_xid` (`xid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

-- the table to store lock data
CREATE TABLE IF NOT EXISTS `lock_table`
(
    `row_key`        VARCHAR(256) NOT NULL,
    `xid`            VARCHAR(128) NOT NULL,
    `transaction_id` BIGINT,
    `branch_id`      BIGINT       NOT NULL,
    `resource_id`    VARCHAR(256),
    `table_name`     VARCHAR(64),
    `pk`             VARCHAR(36),
    `gmt_create`     DATETIME,
    `gmt_modified`   DATETIME,
    PRIMARY KEY (`row_key`),
    KEY `idx_branch_id` (`branch_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

# --------------------------------------- test ------------------------------------
CREATE database if NOT EXISTS `seata_a` default character set utf8mb4 collate utf8mb4_unicode_ci;
use `seata_a`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `a` (
     `a_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
     `created_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
     `updated_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
     `deleted_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
     PRIMARY KEY (`a_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='a表';

CREATE TABLE `undo_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `branch_id` bigint(20) NOT NULL,
    `xid` varchar(128) NOT NULL,
    `context` varchar(128) NOT NULL,
    `rollback_info` longblob NOT NULL,
    `log_status` int(11) NOT NULL,
    `log_created` datetime NOT NULL,
    `log_modified` datetime NOT NULL,
    `ext` varchar(100) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_unionkey` (`xid`,`branch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `branch_transaction` (
    `sysno` bigint(20) NOT NULL AUTO_INCREMENT,
    `xid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
    `branch_id` bigint(20) NOT NULL,
    `args_json` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `state` tinyint(4) DEFAULT NULL COMMENT '1，初始化；2，已提交；3，已回滚',
    `gmt_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`sysno`) USING BTREE,
    UNIQUE KEY `xid` (`xid`,`branch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='事务记录表';

# --------------------------------------- test ------------------------------------
CREATE database if NOT EXISTS `seata_b` default character set utf8mb4 collate utf8mb4_unicode_ci;
use `seata_b`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `b` (
    `b_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
    `deleted_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`b_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='b表';

CREATE TABLE `undo_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `branch_id` bigint(20) NOT NULL,
    `xid` varchar(128) NOT NULL,
    `context` varchar(128) NOT NULL,
    `rollback_info` longblob NOT NULL,
    `log_status` int(11) NOT NULL,
    `log_created` datetime NOT NULL,
    `log_modified` datetime NOT NULL,
    `ext` varchar(100) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_unionkey` (`xid`,`branch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `branch_transaction` (
    `sysno` bigint(20) NOT NULL AUTO_INCREMENT,
    `xid` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
    `branch_id` bigint(20) NOT NULL,
    `args_json` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `state` tinyint(4) DEFAULT NULL COMMENT '1，初始化；2，已提交；3，已回滚',
    `gmt_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`sysno`) USING BTREE,
    UNIQUE KEY `xid` (`xid`,`branch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='事务记录表';