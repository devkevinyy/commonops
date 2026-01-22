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

CREATE DATABASE IF NOT EXISTS commonops;

USE commonops;

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

DROP TABLE IF EXISTS `diff_caches`;

CREATE TABLE `diff_caches` (
  `type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `instance_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
  `tags` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ecs_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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

DROP TABLE IF EXISTS `permissions`;

CREATE TABLE `permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `can_delete` tinyint(4) DEFAULT '1',
  `auth_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `permissions` WRITE;

INSERT INTO `permissions` (`id`, `name`, `description`, `url_path`, `can_delete`, `auth_type`)
VALUES
	(1,'权限-权限链接','权限-权限链接','/admin/permission/permissions',0,'菜单'),
	(2,'权限-角色管理','权限-角色管理','/admin/permission/roles',0,'菜单'),
	(3,'云服务器','云服务器','/admin/cloud_resource/cloud_server',0,'菜单'),
	(4,'云数据库','云数据库','/admin/cloud_resource/cloud_rds',0,'菜单'),
	(5,'KV-Store','KV-Store','/admin/cloud_resource/cloud_kv',0,'菜单'),
	(6,'负载均衡','负载均衡','/admin/cloud_resource/cloud_slb',0,'菜单'),
	(7,'云账号','云账号','/admin/cloud_resource/cloud_account',0,'菜单'),
	(8,'权限-用户管理','权限-用户管理','/admin/permission/users',0,'菜单'),
	(9,'提交工单','提交工单','/admin/task/deploy_project',0,'菜单'),
	(10,'工单列表','工单列表','/admin/task/jobs',0,'菜单'),
	(11,'其它资源','其它资源','/admin/cloud_resource/cloud_other',0,'菜单'),
	(12,'同步阿里云数据','同步阿里云数据','/admin/data/syncAliyun',0,'菜单'),
	(14,'用户反馈','用户反馈','/admin/system/user_feedback',0,'菜单'),
	(15,'kubernetes管理','kubernetes管理入口','/admin/k8s_cluster/info',0,'菜单'),
	(17,'dms权限管理','dms权限管理','/admin/dms/auth_manage',0,'菜单'),
	(18,'dms实例管理','dms实例管理','/admin/dms/instance_manage',0,'菜单'),
	(19,'dms数据操作','dms数据操作','/admin/dms/data_manage',0,'菜单'),
	(20,'DNS解析入口','DNS解析入口','/admin/dns/domain_manage',0,'菜单'),
	(21,'持续集成入口','持续集成入口','/admin/cicd/ci',0,'菜单'),
	(22,'持续部署入口','持续部署入口','/admin/cicd/cd',0,'菜单'),
	(23,'发布记录管理','发布记录管理','/admin/cicd/cd_record',0,'菜单'),
	(24,'获取服务器数据','获取服务器数据','GET:/cloud/servers',0,'操作'),
	(26,'获取权限列表数据','获取权限列表数据','GET:/permissions/list',0,'操作'),
	(27,'获取角色列表数据','获取角色列表数据','GET:/roles/list',0,'操作'),
	(28,'获取角色链接数据','获取角色链接数据','GET:/roles/authLink',0,'操作'),
	(29,'获取角色下的用户数据','获取角色下的用户数据','GET:/roles/users',0,'操作'),
	(30,'获取角色下的资源数据','获取角色下的资源数据','GET:/roles/resources',0,'操作'),
	(31,'获取用户列表数据','获取用户列表数据','GET:/user/list',0,'操作'),
	(32,'获取云账号数据','获取云账号数据','GET:/cloud/accounts',0,'操作'),
	(33,'新增服务器信息','新增服务器信息','POST:/cloud/servers',0,'操作'),
	(34,'修改服务器信息','修改服务器信息','PUT:/cloud/servers',0,'操作'),
	(35,'获取RDS数据','获取RDS数据','GET:/cloud/rds',0,'操作'),
	(36,'新增RDS数据','新增RDS数据','POST:/cloud/rds',0,'操作'),
	(37,'修改RDS数据','修改RDS数据','PUT:/cloud/rds',0,'操作'),
	(38,'获取RDS详情信息','获取RDS详情信息','GET:/cloud/rds/detail',0,'操作'),
	(39,'删除RDS信息','删除RDS信息','DELETE:/cloud/rds',0,'操作'),
	(40,'获取KV信息','获取KV信息','GET:/cloud/kv',0,'操作'),
	(41,'新增KV信息','新增KV信息','POST:/cloud/kv',0,'操作'),
	(42,'修改KV信息','修改KV信息','PUT:/cloud/kv',0,'操作'),
	(43,'删除KV信息','删除KV信息','DELETE:/cloud/kv',0,'操作'),
	(44,'获取KV详情信息','获取KV详情信息','GET:/cloud/kv/detail',0,'操作'),
	(45,'获取SLB信息','获取SLB信息','GET:/cloud/slb',0,'操作'),
	(46,'删除SLB信息','删除SLB信息','DELETE:/cloud/slb',0,'操作'),
	(47,'获取云资源其它信息','获取云资源其它信息','GET:/cloud/other',0,'操作'),
	(48,'新增云资源其它信息','新增云资源其它信息','POST:/cloud/other',0,'操作'),
	(49,'修改云资源其它信息','修改云资源其它信息','PUT:/cloud/other',0,'操作'),
	(50,'删除云资源其它信息','删除云资源其它信息','DELETE:/cloud/other',0,'操作'),
	(51,'新增云账号信息','新增云账号信息','POST:/cloud/accounts',0,'操作'),
	(52,'修改云账号信息','修改云账号信息','PUT:/cloud/accounts',0,'操作'),
	(53,'删除云账号信息','删除云账号信息','DELETE:/cloud/accounts',0,'操作'),
	(54,'获取工单列表信息','获取工单列表信息','GET:/dailyJob/list',0,'操作'),
	(55,'获取工单详情信息','获取工单详情信息','GET:/dailyJob/info',0,'操作'),
	(56,'获取DMS实例信息','获取DMS实例信息','GET:/dms/instances',0,'操作'),
	(57,'获取DMS用户权限信息','获取DMS用户权限信息','GET:/dms/authData',0,'操作'),
	(58,'获取DMS所有实例信息','获取DMS所有实例信息','GET:/dms/instances/all',0,'操作'),
	(59,'获取DMS数据库信息','获取DMS数据库信息','GET:/dms/databaseData',0,'操作'),
	(60,'新增DMS用户权限','新增DMS用户权限','POST:/dms/auth',0,'操作'),
	(61,'删除DMS用户权限','删除DMS用户权限','DELETE:/dms/auth',0,'操作'),
	(62,'获取用户可见实例数据','获取用户可见实例数据','GET:/dms/userInstanceData',0,'操作'),
	(63,'获取DMS执行记录','获取DMS执行记录','GET:/dms/userLog',0,'操作'),
	(64,'允许用户提交SQL','允许用户提交SQL','POST:/dms/userExecSQL',0,'操作'),
	(65,'修改DMS实例','修改DMS实例','PUT:/dms/instance',0,'操作'),
	(66,'删除DMS实例','删除DMS实例','DELETE:/dms/instance',0,'操作'),
	(67,'创建DMS实例','创建DMS实例','POST:/dms/instance',0,'操作'),
	(68,'获取域名列表','获取域名列表','GET:/dns/domainList',0,'操作'),
	(69,'新增域名','新增域名','POST:/dns/domain',0,'操作'),
	(70,'获取域名解析历史记录','获取域名解析历史记录','GET:/dns/domainHistoryList',0,'操作'),
	(71,'获取域名解析列表','获取域名解析列表','GET:/dns/domainRecordsList',0,'操作'),
	(72,'新增域名解析记录','新增域名解析记录','POST:/dns/domainRecord',0,'操作'),
	(73,'修改域名解析记录','修改域名解析记录','POST:/dns/domainRecordUpdate',0,'操作'),
	(74,'删除域名解析记录','删除域名解析记录','DELETE:/dns/domainRecord',0,'操作'),
	(75,'修改域名解析状态','修改域名解析状态','POST:/dns/domainRecordStatus',0,'操作'),
	(76,'获取任务列表','获取任务列表','GET:/ci/jobList',0,'操作'),
	(77,'获取访问凭证列表','获取访问凭证列表','GET:/ci/credentials/list',0,'操作'),
	(78,'获取任务详情','获取任务详情','GET:/ci/job',0,'操作'),
	(79,'创建构建任务','创建构建任务','POST:/ci/job',0,'操作'),
	(80,'删除构建任务','删除构建任务','DELETE:/ci/job',0,'操作'),
	(81,'修改构建任务','修改构建任务','PUT:/ci/job',0,'操作'),
	(82,'获取构建任务列表','获取构建任务列表','GET:/ci/buildList',0,'操作'),
	(83,'获取构建日志','获取构建日志','GET:/ci/buildLog',0,'操作'),
	(84,'获取构建阶段信息','获取构建阶段信息','GET:/ci/build/stages',0,'操作'),
	(85,'获取构建制品信息','获取构建制品信息','GET:/ci/build/archiveArtifactsInfo',0,'操作'),
	(86,'获取部署模板列表','获取部署模板列表','GET:/cd/processTemplateList',0,'操作'),
	(87,'获取K8S集群','获取K8S集群','GET:/kubernetes/cluster',0,'操作'),
	(88,'获取K8S命名空间','获取K8S命名空间','GET:/kubernetes/namespaces',0,'操作'),
	(89,'使用部署模板部署项目','使用部署模板部署项目','POST:/cd/processLog',0,'操作'),
	(90,'创建构建任务','创建构建任务','POST:/ci/build',0,'操作'),
	(91,'添加K8S集群','添加K8S集群','POST:/kubernetes/cluster',0,'操作'),
	(93,'删除K8S集群','删除K8S集群','DELETE:/kubernetes/cluster',0,'操作'),
	(95,'添加K8S命名空间','获取K8S命名空间','POST:/kubernetes/namespaces',0,'操作'),
	(96,'删除K8S命名空间','删除K8S命名空间','DELETE:/kubernetes/namespaces',0,'操作'),
	(97,'获取K8S Deployments','获取K8S Deployments','GET:/kubernetes/deployments',0,'操作'),
	(98,'K8S deployment重启','K8S deployment重启','PUT:/kubernetes/deployment/restart',0,'操作'),
	(99,'获取K8S ReplicationControllers','获取K8S ReplicationControllers','GET:/kubernetes/replication_controllers',0,'操作'),
	(100,'获取K8S ReplicaSets','获取K8S ReplicaSets','GET:/kubernetes/replica_sets',0,'操作'),
	(101,'获取K8S Pods','获取K8S Pods','GET:/kubernetes/pods',0,'操作'),
	(102,'获取K8S Pod Log','获取K8S Pod Log','GET:/kubernetes/pod/log',0,'操作'),
	(103,'获取K8S Nodes','获取K8S Nodes','GET:/kubernetes/nodes',0,'操作'),
	(104,'获取K8S Services','获取K8S Services','GET:/kubernetes/services',0,'操作'),
	(105,'获取K8S Ingress','获取K8S Ingress','GET:/kubernetes/ingress',0,'操作'),
	(106,'获取K8S ConfigMap','获取K8S ConfigMap','GET:/kubernetes/config_dict',0,'操作'),
	(107,'获取K8S Secret','获取K8S Secret','GET:/kubernetes/secret_dict',0,'操作'),
	(108,'添加 K8S YAML资源','添加 K8S YAML资源','POST:/kubernetes/yaml_resource',0,'操作'),
	(109,'修改 K8S YAML资源','修改 K8S YAML资源','PUT:/kubernetes/yaml_resource',0,'操作'),
	(110,'获取 K8S YAML资源','获取 K8S YAML资源','GET:/kubernetes/yaml',0,'操作'),
	(111,'K8S 扩缩容','K8S 扩缩容','PUT:/kubernetes/scale',0,'操作'),
	(112,'删除 K8S 资源','删除 K8S 资源','DELETE:/kubernetes/resource',0,'操作'),
	(113,'删除 K8S ConfigMap','删除 K8S ConfigMap','DELETE:/kubernetes/config_map',0,'操作'),
	(114,'删除 K8S Secret','删除 K8S Secret','DELETE:/kubernetes/secret',0,'操作'),
	(115,'获取 K8S Node监控','获取 K8S Node监控','GET:/kubernetes/metrics/node',0,'操作'),
	(116,'获取 K8S Pod监控','获取 K8S Pod监控','GET:/kubernetes/metrics/pod',0,'操作'),
	(117,'WebSSH:K8S-Pod','WebSSH:K8S-Pod','WebSSH:K8S-Pod',0,'操作'),
	(118,'获取用户反馈列表','获取用户反馈列表','GET:/user/feedback',0,'操作'),
	(119,'配置中心Nacos菜单','配置中心Nacos菜单','/admin/config_center/nacos',0,'菜单'),
	(120,'获取Nacos集群列表','获取Nacos集群列表','GET:/configCenter/nacos/list',0,'操作'),
	(121,'新增Nacos集群','新增Nacos集群','POST:/configCenter/nacos',0,'操作'),
	(122,'获取Nacos命名空间列表','获取Nacos命名空间列表','GET:/configCenter/nacos/namespaces',0,'操作'),
	(123,'获取Nacos命名空间下的配置','获取Nacos命名空间下的配置','GET:/configCenter/nacos/configs',0,'操作'),
	(124,'创建Nacos配置','创建Nacos配置','POST:/configCenter/nacos/config',0,'操作'),
	(125,'修改Nacos配置','修改Nacos配置','PUT:/configCenter/nacos/config',0,'操作'),
	(126,'删除Nacos配置','删除Nacos配置','DELETE:/configCenter/nacos/config',0,'操作'),
	(127,'复制Nacos配置','复制Nacos配置','POST:/configCenter/nacos/config/copy',0,'操作'),
	(128,'获取Nacos配置详情','获取Nacos配置详情','GET:/configCenter/nacos/config',0,'操作'),
	(129,'获取Nacos所有配置','获取Nacos所有配置','GET:/configCenter/nacos/configs/all',0,'操作'),
	(130,'同步Nacos静态配置','同步Nacos静态配置','POST:/configCenter/nacos/config/sync',0,'操作'),
	(131,'获取ECS监控信息','获取ECS监控信息','GET:/cloud/monitor/ecs',0,'操作'),
	(132,'获取ECS实例详情','获取ECS实例详情','GET:/cloud/server',0,'操作'),
	(133,'删除服务器信息','删除服务器信息','DELETE:/cloud/servers',0,'操作'),
	(134,'批量服务器执行命令','批量服务器执行命令','POST:/cloud/servers/batch/ssh',0,'操作'),
	(135,'获取全部服务器树状数据','获取全部服务器树状数据','GET:/cloud/servers/treedata',0,'操作'),
	(136,'配置模板管理','配置模板管理','/admin/config_center/config_template',0,'菜单'),
	(137,'获取配置模板数据','获取配置模板数据','GET:/configCenter/configTemplates',0,'操作'),
	(138,'添加配置模板','添加配置模板','POST:/configCenter/configTemplate',0,'操作'),
	(139,'修改配置模板','修改配置模板','PUT:/configCenter/configTemplate',0,'操作'),
	(140,'删除配置模板','删除配置模板','DELETE:/configCenter/configTemplate',0,'操作'),
	(141,'获取所有模板数据','获取所有模板数据','GET:/configCenter/configTemplates/all',0,'操作'),
	(142,'添加角色','添加角色','POST:/roles/addRole',0,'操作'),
	(143,'修改角色','修改角色','PUT:/roles/updateRole',0,'操作'),
	(144,'删除角色','删除角色','DELETE:/roles/deleteRole',0,'操作'),
	(145,'添加角色用户','添加角色用户','POST:/roles/users',0,'操作'),
	(146,'添加角色资源','添加角色资源','POST:/roles/resources',0,'操作'),
	(147,'添加用户权限链接','添加用户权限链接','POST:/roles/authLink',0,'操作'),
	(148,'添加角色链接','添加角色链接','POST:/roles/authLinks',0,'操作'),
	(149,'删除权限链接','删除权限链接','DELETE:/roles/authLink',0,'操作'),
	(150,'用户token刷新','用户token刷新','GET:/user/tokenRefresh',0,'操作'),
	(151,'用户修改密码','用户修改密码','POST:/user/updatePassword',0,'操作'),
	(152,'创建新用户','创建新用户','POST:/user/create',0,'操作'),
	(153,'用户状态激活','用户状态激活','POST:/user/active',0,'操作'),
	(154,'获取用户权限','获取用户权限','GET:/user/permissions',0,'操作'),
	(155,'创建用户反馈','创建用户反馈','POST:/user/feeback',0,'操作'),
	(156,'获取RDS监控信息','获取RDS监控信息','GET:/cloud/monitor/rds',0,'操作'),
	(157,'获取KV监控信息','获取KV监控信息','GET:/cloud/monitor/kv',0,'操作'),
	(158,'获取SLB监控信息','获取SLB监控信息','GET:/cloud/monitor/slb',0,'操作'),
	(159,'创建工单','创建工单','POST:/dailyJob/',0,'操作'),
	(160,'修改工单','修改工单','PUT:/dailyJob/',0,'操作'),
	(161,'修改工单执行用户','修改工单执行用户','PUT:/dailyJob/executorUser',0,'操作'),
	(162,'同步阿里云ecs数据','同步阿里云ecs数据','GET:/data/syncAliyunEcs',0,'操作'),
	(163,'同步阿里云rds数据','同步阿里云rds数据','GET:/data/syncAliyunRds',0,'操作'),
	(164,'同步阿里云kv数据','同步阿里云kv数据','GET:/data/syncAliyunKv',0,'操作'),
	(165,'同步阿里云slb数据','同步阿里云slb数据','GET:/data/syncAliyunSlb',0,'操作'),
	(166,'同步阿里云统计数据','同步阿里云统计数据','GET:/data/syncAliyunStatisData',0,'操作'),
	(167,'获取DMS实例数据','获取DMS实例数据','GET:/dms/instanceData',0,'操作'),
	(168,'删除DMS数据库数据','删除DMS数据库数据','DELETE:/dms/databaseData',0,'操作'),
	(169,'新增DMS数据库数据','新增DMS数据库数据','POST:/dms/databaseData',0,'操作'),
	(170,'获取DMS用户数据库数据','获取DMS用户数据库数据','GET:/dms/userDatabaseData',0,'操作'),
	(171,'获取CI Build信息','获取CI Build信息','GET:/ci/buildInfo',0,'操作'),
	(172,'获取CI BuildStageLog','获取CI BuildStageLog','GET:/ci/build/state/log',0,'操作'),
	(173,'删除CI 构建','删除CI 构建','DELETE:/ci/build',0,'操作'),
	(174,'CI创建凭证','CI创建凭证','POST:/ci/credential',0,'操作'),
	(175,'新建CD流程模板','新建CD流程模板','POST:/cd/processTemplate',0,'操作'),
	(176,'获取CD流程日志','获取CD流程日志','GET:/cd/processLog',0,'操作'),
    (177,'更新用户权限','更新用户权限','PUT:/user/active',0,'操作'),
    (178,'获取权限链接-新','获取权限链接-新','GET:/permissions/authLink',0,'操作'),
    (179,'添加权限链接-新','添加权限链接-新','POST:/permissions/authLink',0,'操作'),
    (180,'更新权限链接-新','更新权限链接-新','PUT:/permissions/authLink',0,'操作'),
    (181,'删除权限链接-新','删除权限链接-新','DELETE:/permissions/authLink',0,'操作'),
    (182,'Ansible菜单栏','Ansible菜单栏','/admin/batch/ansible',0,'菜单'),
    (183,'获取节点树','获取节点树','GET:/batch/allNodeTree',0,'操作'),
    (184,'批量命令菜单栏','批量命令菜单栏','/admin/batch/cmds',0,'菜单');

UNLOCK TABLES;

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

DROP TABLE IF EXISTS `role_permissions`;

CREATE TABLE `role_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) DEFAULT NULL,
  `permission_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `role_permissions` WRITE;

INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`)
VALUES
	(178,1,23),
	(179,1,21),
	(180,1,22),
	(181,1,20),
	(182,1,19),
	(183,1,17),
	(184,1,18),
	(185,1,16),
	(186,1,15),
	(187,1,14),
	(188,1,13),
	(189,1,12),
	(190,1,11),
	(191,1,3),
	(192,1,4),
	(193,1,5),
	(194,1,6),
	(195,1,7),
	(196,1,8),
	(197,1,9),
	(198,1,10),
	(199,1,1),
	(200,1,2),
	(202,1,24),
	(203,1,26),
	(204,1,27),
	(205,1,28),
	(206,1,29),
	(207,1,30),
	(208,1,31),
	(209,1,32),
	(210,1,33),
	(211,1,34),
	(212,1,35),
	(213,1,36),
	(214,1,37),
	(215,1,38),
	(216,1,39),
	(217,1,40),
	(218,1,41),
	(219,1,42),
	(220,1,43),
	(221,1,44),
	(222,1,45),
	(223,1,46),
	(224,1,47),
	(225,1,48),
	(226,1,49),
	(227,1,50),
	(228,1,51),
	(229,1,52),
	(230,1,53),
	(231,1,54),
	(232,1,55),
	(233,1,56),
	(234,1,57),
	(235,1,58),
	(236,1,59),
	(237,1,60),
	(238,1,61),
	(239,1,62),
	(240,1,63),
	(241,1,64),
	(242,1,65),
	(243,1,66),
	(244,1,67),
	(245,1,68),
	(246,1,69),
	(247,1,70),
	(248,1,71),
	(249,1,72),
	(250,1,73),
	(251,1,74),
	(252,1,75),
	(253,1,76),
	(254,1,77),
	(255,1,78),
	(256,1,79),
	(257,1,80),
	(258,1,81),
	(259,1,82),
	(260,1,83),
	(261,1,84),
	(262,1,85),
	(263,1,86),
	(264,1,87),
	(265,1,88),
	(266,1,89),
	(267,1,90),
	(268,1,91),
	(270,1,93),
	(272,1,109),
	(273,1,95),
	(274,1,96),
	(275,1,97),
	(276,1,98),
	(277,1,99),
	(278,1,100),
	(279,1,101),
	(280,1,102),
	(281,1,103),
	(282,1,104),
	(283,1,105),
	(284,1,106),
	(285,1,107),
	(286,1,108),
	(287,1,110),
	(288,1,111),
	(289,1,112),
	(290,1,113),
	(291,1,114),
	(292,1,115),
	(293,1,116),
	(294,1,117),
	(295,1,118),
	(296,1,119),
	(297,1,120),
	(298,1,121),
	(299,1,122),
	(300,1,123),
	(301,1,124),
	(302,1,125),
	(303,1,126),
	(304,1,127),
	(305,1,128),
	(306,1,129),
	(307,1,130),
	(308,1,131),
	(309,1,132),
	(310,1,133),
	(311,1,134),
	(312,1,135),
	(313,1,136),
	(314,1,137),
	(315,1,138),
	(316,1,139),
	(317,1,140),
	(318,1,141),
	(319,1,142),
	(320,1,143),
	(321,1,144),
	(322,1,145),
	(323,1,146),
	(324,1,147),
	(325,1,148),
	(326,1,149),
	(327,1,150),
	(328,1,151),
	(329,1,152),
	(330,1,153),
	(331,1,154),
	(332,1,155),
	(333,1,156),
	(334,1,157),
	(335,1,158),
	(336,1,159),
	(337,1,160),
	(338,1,161),
	(339,1,162),
	(340,1,163),
	(341,1,164),
	(342,1,165),
	(343,1,166),
	(344,1,167),
	(345,1,168),
	(346,1,169),
	(347,1,170),
	(348,1,171),
	(349,1,172),
	(350,1,173),
	(351,1,174),
	(352,1,175),
	(353,1,176),
    (354,1,177),
    (355,1,178),
    (356,1,179),
    (357,1,180),
    (358,1,181),
    (359,1,182),
    (360,1,183),
    (361,1,184);

UNLOCK TABLES;

DROP TABLE IF EXISTS `role_resources`;

CREATE TABLE `role_resources` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) DEFAULT NULL,
  `resource_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `resource_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `roles` WRITE;

INSERT INTO `roles` (`id`, `name`, `description`)
VALUES
	(1,'超级管理员','超级管理员');

UNLOCK TABLES;

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

DROP TABLE IF EXISTS `user_roles`;

CREATE TABLE `user_roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

LOCK TABLES `user_roles` WRITE;

INSERT INTO `user_roles` (`id`, `user_id`, `role_id`)
VALUES
	(2,1,1);

UNLOCK TABLES;

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

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `union_id`, `mobile`, `is_leader`, `user_name`, `password`, `emp_id`, `department_id`, `active`, `position`, `email`, `first_in`)
VALUES
	(1,'2019-06-05 16:08:58','2019-09-05 23:52:11',NULL,'','','',0,'admin','24bf71e48e2fde8675fe86c320ce559f','1',NULL,1,'系统管理员','admin@ops.com',0);

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

DROP TABLE IF EXISTS `cd_process_log`;

CREATE TABLE `cd_process_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `emp_id` varchar(32) DEFAULT NULL,
  `template_id` int(11) DEFAULT NULL,
  `image_name` varchar(256) DEFAULT NULL,
  `success` int(11) DEFAULT NULL,
  `result` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cd_process_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `cd_process_template`;

CREATE TABLE `cd_process_template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `emp_id` varchar(32) DEFAULT NULL,
  `template_name` varchar(512) DEFAULT NULL,
  `cluster_id` varchar(128) DEFAULT NULL,
  `namespace` varchar(32) DEFAULT NULL,
  `deploy_yaml` text,
  `service_yaml` text,
  `configmap_yaml` text,
  `ingress_yaml` text,
  `job_name` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cd_process_template_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `nacos`;

CREATE TABLE `nacos` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `data_status` int(11) NOT NULL DEFAULT '1',
  `end_point` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `alias` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `nacos_config`;

CREATE TABLE `nacos_config` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `cluster_id` int(11) DEFAULT NULL,
  `namespace` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `data_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `config_group` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `template_id` int(11) DEFAULT NULL,
  `fill_data` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `data_status` tinyint(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `config_template`;

CREATE TABLE `config_template` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `data_status` tinyint(4) NOT NULL DEFAULT '1',
  `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `config_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `fill_field` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
