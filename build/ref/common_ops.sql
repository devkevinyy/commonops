# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.24)
# Database: common_ops
# Generation Time: 2019-11-19 03:44:51 +0000
# ************************************************************

SET NAMES utf8mb4;

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table cloud_account
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cloud_account`;

CREATE TABLE `cloud_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cloud_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `passwd` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `secret` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table daily_jobs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `daily_jobs`;

CREATE TABLE `daily_jobs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `job_name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `job_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `important_degree` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `open_deploy_auto_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `task_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `creator_user_id` int(11) DEFAULT NULL,
  `creator_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `executor_user_id` int(11) DEFAULT NULL,
  `executor_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL,
  `accept_time` timestamp NULL DEFAULT NULL,
  `end_time` timestamp NULL DEFAULT NULL,
  `refuse_reason` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table diff_caches
# ------------------------------------------------------------

DROP TABLE IF EXISTS `diff_caches`;

CREATE TABLE `diff_caches` (
  `type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table dms_approve
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dms_approve`;

CREATE TABLE `dms_approve` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `emp_id` varchar(32) NOT NULL,
  `username` varchar(256) NOT NULL,
  `instance_type` varchar(256) NOT NULL,
  `instance_id` varchar(512) NOT NULL,
  `instance_name` varchar(512) NOT NULL,
  `database_id` varchar(256) NOT NULL,
  `database_name` varchar(256) NOT NULL,
  `sql_content` text NOT NULL,
  `sql_description` varchar(2048) DEFAULT NULL,
  `create_time` varchar(256) NOT NULL,
  `approve_id` varchar(256) NOT NULL DEFAULT '',
  `approve_content` varchar(128) NOT NULL DEFAULT '审批中',
  `approve_status` tinyint(4) NOT NULL DEFAULT '0',
  `log_id` int(11) NOT NULL,
  `has_executed` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_dms_approve_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table dms_auth
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dms_auth`;

CREATE TABLE `dms_auth` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `emp_id` int(11) NOT NULL,
  `username` varchar(256) NOT NULL,
  `auth_type` tinyint(4) NOT NULL DEFAULT '1',
  `instance_type` varchar(256) NOT NULL,
  `instance_id` varchar(512) NOT NULL,
  `instance_name` varchar(512) NOT NULL,
  `database_id` varchar(256) NOT NULL,
  `database_name` varchar(256) NOT NULL,
  `oper_type` tinyint(4) NOT NULL DEFAULT '1',
  `valid_time` varchar(256) NOT NULL,
  `oper_count` int(11) NOT NULL DEFAULT '5',
  `allow_tables` varchar(2048) NOT NULL,
  `approve_emp_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_dms_auth_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table dms_database
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dms_database`;

CREATE TABLE `dms_database` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `instance_id` varchar(128) DEFAULT NULL,
  `database_id` varchar(128) DEFAULT NULL,
  `schema_name` varchar(128) DEFAULT NULL,
  `state` varchar(32) DEFAULT NULL,
  `db_type` varchar(32) DEFAULT NULL,
  `host` varchar(256) DEFAULT NULL,
  `port` smallint(6) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `database_id_idx` (`database_id`),
  KEY `idx_dms_database_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table dms_instance
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dms_instance`;

CREATE TABLE `dms_instance` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `port` smallint(6) NOT NULL,
  `instance_type` varchar(32) NOT NULL,
  `host` varchar(128) NOT NULL,
  `state` varchar(32) NOT NULL,
  `instance_id` varchar(128) NOT NULL,
  `instance_alias` varchar(128) NOT NULL,
  `oper_user` varchar(128) NOT NULL,
  `oper_pwd` varchar(256) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `instance_id_idx` (`instance_id`),
  KEY `idx_dms_instance_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table dms_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dms_log`;

CREATE TABLE `dms_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `emp_id` varchar(32) DEFAULT NULL,
  `username` varchar(32) DEFAULT NULL,
  `database_id` varchar(128) DEFAULT NULL,
  `database_name` varchar(32) DEFAULT NULL,
  `start_time` varchar(32) DEFAULT NULL,
  `sql_content` text,
  `exec_status` tinyint(4) DEFAULT NULL,
  `duration` int(11) DEFAULT NULL,
  `effect_rows` int(11) DEFAULT NULL,
  `result` text,
  `exception_output` text,
  `has_executed` tinyint(4) NOT NULL DEFAULT '0',
  `rollback_table_name` varchar(1024) DEFAULT NULL,
  `has_rollback` tinyint(4) NOT NULL DEFAULT '0',
  `rollback_time` varchar(32) NOT NULL,
  `sql_type` varchar(32) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_dms_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table dms_table
# ------------------------------------------------------------

DROP TABLE IF EXISTS `dms_table`;

CREATE TABLE `dms_table` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `table_id` varchar(128) DEFAULT NULL,
  `database_id` varchar(128) DEFAULT NULL,
  `table_name` varchar(32) DEFAULT NULL,
  `table_schema_name` varchar(128) DEFAULT NULL,
  `engine` varchar(128) DEFAULT NULL,
  `encoding` varchar(32) DEFAULT NULL,
  `table_type` varchar(256) DEFAULT NULL,
  `num_rows` varchar(32) DEFAULT NULL,
  `store_capacity` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `table_id_idx` (`table_id`),
  KEY `idx_dms_table_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table ecs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ecs`;

CREATE TABLE `ecs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `cloud_account_id` int(11) DEFAULT NULL,
  `image_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_network_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `local_storage_amount` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_charge_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cluster_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `start_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `internet_charge_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `internet_max_bandwidth_in` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `host_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cpu` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `os_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `os_name_en` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `serial_number` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `internet_max_bandwidth_out` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `resource_group_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_type_family` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deployment_set_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `expired_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `os_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `memory` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `creation_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `local_storage_capacity` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `inner_ip_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `public_ip_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `private_ip_address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ssh_port` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `ssh_user` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `ssh_pwd` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ecs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table kv
# ------------------------------------------------------------

DROP TABLE IF EXISTS `kv`;

CREATE TABLE `kv` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `cloud_account_id` int(11) DEFAULT NULL,
  `instance_class` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `package_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `charge_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `connection_domain` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `search_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `has_renew_change_order` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `destroy_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `private_ip` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bandwidth` int(11) DEFAULT NULL,
  `network_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `vpc_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `node_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `connections` int(11) DEFAULT NULL,
  `architecture_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `replacate_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `engine_version` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `capacity` int(11) DEFAULT NULL,
  `v_switch_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `port` int(11) DEFAULT NULL,
  `zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `end_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `qps` int(11) DEFAULT NULL,
  `user_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `config` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_rds` tinyint(1) DEFAULT NULL,
  `connection_mode` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_kv_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table other_res
# ------------------------------------------------------------

DROP TABLE IF EXISTS `other_res`;

CREATE TABLE `other_res` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `cloud_account_id` int(11) DEFAULT NULL,
  `res_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `connections` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `engine` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cpu` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `disk` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bandwidth` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `memory` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `renew_status` int(11) NOT NULL DEFAULT '1',
  `create_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `expired_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_other_res_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table permissions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `permissions`;

CREATE TABLE `permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url_path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;

INSERT INTO `permissions` (`id`, `name`, `description`, `url_path`)
VALUES
	(1,'权限-权限链接','权限-权限链接','/admin/permission/permissions'),
	(2,'权限-角色管理','权限-角色管理','/admin/permission/roles'),
	(3,'云服务器','云服务器','/admin/cloud_resource/cloud_server'),
	(4,'云数据库','云数据库','/admin/cloud_resource/cloud_rds'),
	(5,'KV-Store','KV-Store','/admin/cloud_resource/cloud_kv'),
	(6,'负载均衡','负载均衡','/admin/cloud_resource/cloud_slb'),
	(7,'云账号','云账号','/admin/cloud_resource/cloud_account'),
	(8,'权限-用户管理','权限-用户管理','/admin/permission/users'),
	(9,'提交工单','提交工单','/admin/task/deploy_project'),
	(10,'工单列表','工单列表','/admin/task/jobs'),
	(11,'其它资源','其它资源','/admin/cloud_resource/cloud_other'),
	(12,'同步阿里云数据','同步阿里云数据','/admin/data/syncAliyun'),
	(13,'CI&CD','CI&CD访问入口','/admin/cicd/jobs'),
	(14,'用户反馈','用户反馈','/admin/system/user_feedback'),
	(15,'kubernetes管理','kubernetes管理入口','/admin/k8s_cluster/info'),
	(17,'dms权限管理','dms权限管理','/admin/dms/auth_manage'),
	(18,'dms实例管理','dms实例管理','/admin/dms/instance_manage'),
	(19,'dms数据操作','dms数据操作','/admin/dms/data_manage');

/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table rds
# ------------------------------------------------------------

DROP TABLE IF EXISTS `rds`;

CREATE TABLE `rds` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `cloud_account_id` int(11) DEFAULT NULL,
  `ins_id` int(11) DEFAULT NULL,
  `db_instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `expire_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `destroy_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `engine` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_net_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `connection_mode` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `lock_mode` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `category` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_storage_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_class` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_network_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `vpc_cloud_instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `lock_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mutri_o_rsignle` tinyint(1) DEFAULT NULL,
  `create_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `engine_version` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `guard_db_instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `temp_db_instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `master_instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `vpc_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v_switch_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `replicate_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `resource_group_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `connection_string` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `port` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `db_instance_memory` int(11) DEFAULT NULL,
  `db_instance_storage` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rds_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table role_permissions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `role_permissions`;

CREATE TABLE `role_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) DEFAULT NULL,
  `permission_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `role_permissions` WRITE;
/*!40000 ALTER TABLE `role_permissions` DISABLE KEYS */;

INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`)
VALUES
	(78,1,16),
	(79,1,15),
	(80,1,14),
	(81,1,13),
	(82,1,12),
	(83,1,11),
	(84,1,3),
	(85,1,4),
	(86,1,5),
	(87,1,6),
	(88,1,7),
	(89,1,8),
	(90,1,9),
	(91,1,10),
	(92,1,1),
	(93,1,2),
  (94,1,17),
  (95,1,18),
  (96,1,19);

/*!40000 ALTER TABLE `role_permissions` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table role_resources
# ------------------------------------------------------------

DROP TABLE IF EXISTS `role_resources`;

CREATE TABLE `role_resources` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) DEFAULT NULL,
  `resource_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `resource_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



# Dump of table roles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;

INSERT INTO `roles` (`id`, `name`, `description`)
VALUES
	(1,'超级管理员','超级管理员');

/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table slb
# ------------------------------------------------------------

DROP TABLE IF EXISTS `slb`;

CREATE TABLE `slb` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `cloud_account_id` int(11) DEFAULT NULL,
  `count` int(11) DEFAULT NULL,
  `slave_zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `load_balancer_status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `v_switch_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `master_zone_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pay_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region_id_alias` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `load_balancer_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address_ip_version` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `region_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `resource_group_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `load_balancer_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `internet_charge_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `vpc_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `network_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_slb_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table user_feedback
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_feedback`;

CREATE TABLE `user_feedback` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `user_id` int(11) DEFAULT NULL,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `content` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `score` int(11) NOT NULL DEFAULT '0',
  `create_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# Dump of table user_roles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_roles`;

CREATE TABLE `user_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;

INSERT INTO `user_roles` (`id`, `user_id`, `role_id`)
VALUES
	(2,1,1);

/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `union_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mobile` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `is_leader` tinyint(1) DEFAULT NULL,
  `user_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `emp_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `department_id` int(11) DEFAULT NULL,
  `active` tinyint(1) DEFAULT NULL,
  `position` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `first_in` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `union_id`, `mobile`, `is_leader`, `user_name`, `password`, `emp_id`, `department_id`, `active`, `position`, `email`, `first_in`)
VALUES
	(1,'2019-06-05 16:08:58','2019-09-05 23:52:11',NULL,'','','',0,'admin','24bf71e48e2fde8675fe86c320ce559f','1',NULL,1,'系统管理员','admin@ops.com',0);

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

DROP TABLE IF EXISTS `k8s`;

CREATE TABLE `k8s` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cluster_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `token` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `api_server` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

DROP TABLE IF EXISTS `nodes`;

CREATE TABLE `nodes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cid` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cpu` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `memory` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `storage` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `pods`;

CREATE TABLE `pods` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cid` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `namespace` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `container` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cpu` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `memory` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `storage` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;