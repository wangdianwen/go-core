-- MySQL dump 10.13  Distrib 8.0.16, for osx10.14 (x86_64)
--
-- Host: localhost    Database: core
-- ------------------------------------------------------
-- Server version	8.0.16

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8mb4 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

DROP DATABASE IF EXISTS `b_core`;
CREATE DATABASE IF NOT EXISTS `b_core`;
USE `b_core`;


--
-- Table structure for table `core_apps`
--

DROP TABLE IF EXISTS `core_apps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_apps` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'name of the application',
  `app_key` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'key of the application',
  `app_secret` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'secret of the application',
  `ctime` datetime NOT NULL COMMENT 'create time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  `allow_free_trial` tinyint(4) NOT NULL COMMENT 'dose the app allow a free trial, 0 yes 1 no',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'application status 0: normal, 1: forbidden',
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='all registered applications';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_apps`
--

LOCK TABLES `core_apps` WRITE;
/*!40000 ALTER TABLE `core_apps` DISABLE KEYS */;
INSERT INTO
`core_apps`
VALUES
(1,'Core','CoreKEY','e7dd84d141af90a8f2df885120ffc1e5', NOW(), NOW(),1,0),
(2,'OEPickOut','OEPickOutKEY','0544ed93c22d5e1c869010c4515a7401', NOW(), NOW(),0,0);
/*!40000 ALTER TABLE `core_apps` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_client`
--

DROP TABLE IF EXISTS `core_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_client` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'client name',
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'description of the company',
  `back_end_type` tinyint(3) NOT NULL DEFAULT 0 COMMENT 'backend type: 0 mysql, 1 accredo, 2 advanced',
  `connection` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'connection str of the client',
  `phone` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the phone number of the client',
  `fax` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the fax number of the client',
  `email` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the email address of the client',
  `address_line_1` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the address line 1 of the client',
  `address_line_2` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the address line 2 of the client',
  `city` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the city of the client',
  `country` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the country of the client',
  `postal_code` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the postal code of the client',
  `state` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the state of the client',
  `web` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the web address of the client',
  `gst` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the GST number of the client',
  `bank_account` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'the bank account of the client',
  `expire_time` datetime NOT NULL COMMENT 'expired day',
  `ctime` datetime NOT NULL COMMENT 'register time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  `allow_userlog` tinyint(3) NOT NULL DEFAULT '1' COMMENT 'is allow user log, 0: allow, 1 deny',
  `max_users` int(10) NOT NULL DEFAULT '0' COMMENT 'every client allow max user number',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'user status 0: normal, 1: forbidden',
  PRIMARY KEY (`id`),
  KEY `expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='client lists';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `core_menu`
--

DROP TABLE IF EXISTS `core_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_menu` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `appid` bigint(20) NOT NULL COMMENT 'menu belong to which application',
  `paid` bigint(20) NOT NULL DEFAULT '0' COMMENT 'parent menu id',
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'menu name',
  `url` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'menu url',
  `visible` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'is visible;0: visile, 1 hide',
  `ctime` datetime NOT NULL COMMENT 'create time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `searchByName` (`name`),
  KEY `searchByParent` (`paid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='menu lists';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_menu`
--

LOCK TABLES `core_menu` WRITE;
/*!40000 ALTER TABLE `core_menu` DISABLE KEYS */;
INSERT INTO
`core_menu`
VALUES
(1, 1, 0,'admin','/admin',0, NOW(), NOW()),
(2, 1, 1,'users','/admin/users',0,NOW(), NOW()),
(3, 1, 1,'clients','/admin/clients',0,NOW(), NOW()),
(4, 1, 1,'roles','/admin/roles',0,NOW(), NOW()),
(5, 1, 1,'menus','/admin/menus',0,NOW(), NOW()),
(6, 1, 1,'privileges','/admin/privileges',0, NOW(), NOW()),
(7, 1, 1,'apps','/admin/apps',0,NOW(), NOW()),
(8, 2, 0,'oepickout','/oepickout',0,NOW(), NOW()),
(9, 2, 8,'oepickoutlist','/oepickout/oepickoutlist',0,NOW(), NOW()),
(10, 2, 8,'oepickoutdetail','/oepickout/oepickoutdetail',0,NOW(), NOW());
/*!40000 ALTER TABLE `core_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_privilege`
--

DROP TABLE IF EXISTS `core_privilege`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_privilege` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mid` bigint(20) NOT NULL COMMENT 'belong to which menu',
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'privilege name',
  `controller` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'privilege controller name',
  `action` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'privilege action name',
  `ctime` datetime NOT NULL COMMENT 'create time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'user status 0: normal, 1: forbidden',
  PRIMARY KEY (`id`),
  UNIQUE KEY `searchController` (`controller`,`action`),
  UNIQUE KEY `searchByName` (`name`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='core privileges lists';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_privilege`
--

LOCK TABLES `core_privilege` WRITE;
/*!40000 ALTER TABLE `core_privilege` DISABLE KEYS */;
INSERT INTO
`core_privilege`
VALUES
(1, 2,'get user suggestions','UserController','UserSuggestion', NOW(), NOW(),0),
(2, 2,'search user','UserController','SearchUsers', NOW(), NOW(),0),
(3, 2,'create user','UserController','SaveUser',NOW(), NOW(),0),
(4, 2,'edit user','UserController','EditUser',NOW(), NOW(),0),
(5, 2,'forbidden user','UserController','BanUser',NOW(), NOW(),0),
(6, 2,'release user','UserController','ReleaseUser',NOW(), NOW(),0),
(7, 4,'search roles','RoleController','RoleSearch',NOW(), NOW(),0),
(8, 4,'create role','RoleController','SaveRole',NOW(), NOW(),0),
(9, 4,'edit role','RoleController','EditRole',NOW(), NOW(),0),
(10, 4,'ban role','RoleController','BanRoles',NOW(), NOW(),0),
(11, 4,'release role','RoleController','ReleaseRole',NOW(), NOW(),0),
(12, 9,'Get Locations','OEPickOutController','GetLocations',NOW(), NOW(),0),
(13, 9,'Get Periods','OEPickOutController','GetPeriods',NOW(), NOW(),0),
(14, 9,'Get System Period','OEPickOutController','GetSystemPeriod',NOW(), NOW(),0),
(15, 9,'Get Customers','OEPickOutController','GetCustomers',NOW(), NOW(),0),
(16, 9,'Get Sales Orders','OEPickOutController','GetSalesOrders',NOW(), NOW(),0),
(17, 10,'Get Sales Order','OEPickOutController','GetSalesOrder',NOW(), NOW(),0),
(18, 10,'Save Quantity PickOut','OEPickOutController','SaveQuantityPickOut',NOW(), NOW(),0),
(19, 10,'Complete PickOut','OEPickOutController','CompletePickOut',NOW(), NOW(),0);



    /*!40000 ALTER TABLE `core_privilege` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_role`
--

DROP TABLE IF EXISTS `core_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rolename` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'role name',
  `cid` int(10) NOT NULL COMMENT 'company id',
  `ctime` datetime NOT NULL COMMENT 'create time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'user status 0: normal, 1: forbidden',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role` (`rolename`,`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='core users roles';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `core_role_privilege`
--

DROP TABLE IF EXISTS `core_role_privilege`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_role_privilege` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rid` int(10) NOT NULL COMMENT 'role id',
  `pid` int(10) DEFAULT NULL COMMENT 'privilege id',
  `mid` int(10) DEFAULT NULL COMMENT 'menu id',
  `ctime` datetime NOT NULL COMMENT 'register time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_privilege_pk` (`rid`,`mid`,`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='privilege-role relationships';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_role_user`
--

LOCK TABLES `core_role_privilege` WRITE;
/*!40000 ALTER TABLE `core_role_privilege` DISABLE KEYS */;
/*!40000 ALTER TABLE `core_role_privilege` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_role_user`
--

DROP TABLE IF EXISTS `core_role_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_role_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rid` int(10) NOT NULL COMMENT 'role id',
  `uid` int(10) NOT NULL COMMENT 'user id',
  `ctime` datetime NOT NULL COMMENT 'create time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `find_privileges` (`uid`,`rid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='users role relationships';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_role_user`
--

LOCK TABLES `core_role_user` WRITE;
/*!40000 ALTER TABLE `core_role_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `core_role_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_statistics_clients_usage`
--

DROP TABLE IF EXISTS `core_statistics_clients_usage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_statistics_clients_usage` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `cid` bigint(20) NOT NULL COMMENT 'client id',
  `number` bigint(20) NOT NULL DEFAULT '0' COMMENT 'period usage',
  `stime` datetime NOT NULL COMMENT 'start time',
  `etime` datetime NOT NULL COMMENT 'end time',
  `ctime` datetime NOT NULL COMMENT 'create time',
  PRIMARY KEY (`id`),
  KEY `ctime` (`ctime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='clinets usage statistics';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_statistics_clients_usage`
--

LOCK TABLES `core_statistics_clients_usage` WRITE;
/*!40000 ALTER TABLE `core_statistics_clients_usage` DISABLE KEYS */;
/*!40000 ALTER TABLE `core_statistics_clients_usage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_user`
--

DROP TABLE IF EXISTS `core_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uname` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user name for login',
  `passwd` char(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '96e79218965eb72c92a549dd5a330112' COMMENT 'md5 user password',
  `realname` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user real name',
  `email` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user email',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'phone number',
  `avatar` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'user avatar',
  `barcode` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'user barcode, barcode login support',
  `cid` int(10) NOT NULL COMMENT 'which company the user belong',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT 'user status 0: normal, 1: forbidden, 2: free try',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'user type: 0 normal, 1 free trail user, 2 empty user, 3 free trial expire user',
  `email_verify_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'email verify status: 0 verifyed, 1 not verified',
  `source` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'where is the user from: 0 from admin create, 1: from user register',
  `last_login_time` datetime DEFAULT NULL COMMENT 'last login time',
  `last_login_ip` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'last login ip',
  `last_login_token` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'last login token',
  `ctime` datetime NOT NULL COMMENT 'register time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uname` (`uname`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `barcode` (`barcode`),
  KEY `company_users` (`cid`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='users table';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_user`
--

LOCK TABLES `core_user` WRITE;
/*!40000 ALTER TABLE `core_user` DISABLE KEYS */;
INSERT INTO
`core_user` (`id`, `uname`, `passwd`, `realname`, `email`, `phone`, `avatar`, `barcode`, `cid`, `status`, `type`, `email_verify_status`, `source`, `last_login_time`, `last_login_ip`, `last_login_token`, `ctime`, `mtime`)
VALUES
(null, 'admin','21232f297a57a5a743894a0e4a801fc3','administrator','admin@brunton.co.nz','1234567890','','1234',-1,0,0,0,0, NOW(), '127.0.0.1', '', NOW(), NOW());
/*!40000 ALTER TABLE `core_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_user_free_trial_products`
--

DROP TABLE IF EXISTS `core_user_free_trial_products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_user_free_trial_products` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT 'user id',
  `appid` bigint(20) NOT NULL COMMENT 'application id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `find_products` (`uid`,`appid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='user free trial products lists';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_user_free_trial_products`
--

LOCK TABLES `core_user_free_trial_products` WRITE;
/*!40000 ALTER TABLE `core_user_free_trial_products` DISABLE KEYS */;
/*!40000 ALTER TABLE `core_user_free_trial_products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `core_user_log`
--

DROP TABLE IF EXISTS `core_user_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `core_user_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT 'user id',
  `pid` bigint(20) NOT NULL COMMENT 'what kind of privilege they use',
  `data` text COLLATE utf8mb4_unicode_ci COMMENT 'json, input all insert or update data',
  `ctime` datetime NOT NULL COMMENT 'register time',
  PRIMARY KEY (`id`),
  KEY `uid_pid` (`uid`,`pid`,`ctime`),
  KEY `uid_ctime` (`uid`,`ctime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='user log';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `core_user_log`
--

LOCK TABLES `core_user_log` WRITE;
/*!40000 ALTER TABLE `core_user_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `core_user_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `services_captcha`
--

DROP TABLE IF EXISTS `services_captcha`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `services_captcha` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'application name',
  `mark` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'check mark, could be username, user email or even ip address',
  `code` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'generate code',
  `ctime` datetime NOT NULL COMMENT 'create time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  `expire` datetime NOT NULL COMMENT 'expire time',
  `verify` datetime DEFAULT NULL COMMENT 'verify time',
  `try_times` int(10) NOT NULL DEFAULT '0' COMMENT 'try times',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0 not verified, 1: verified, 2: closed, 3: cancel',
  PRIMARY KEY (`id`),
  KEY `find_captcha` (`mark`,`app`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='basic service: captcha service';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `services_captcha`
--

LOCK TABLES `services_captcha` WRITE;
/*!40000 ALTER TABLE `services_captcha` DISABLE KEYS */;
/*!40000 ALTER TABLE `services_captcha` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `services_task`
--

DROP TABLE IF EXISTS `services_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `services_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'application name',
  `type` tinyint(3) NOT NULL DEFAULT '1' COMMENT '1: send email,',
  `data` json NOT NULL COMMENT 'task data',
  `ctime` datetime NOT NULL COMMENT 'create time',
  `mtime` datetime NOT NULL COMMENT 'modify time',
  `try_times` int(10) NOT NULL DEFAULT '0' COMMENT 'try times',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0 execute, 1: finished, 2: closed, 3: cancel',
  `reason` text COLLATE utf8mb4_unicode_ci COMMENT 'fail reason',
  PRIMARY KEY (`id`),
  KEY `find_queue` (`status`,`app`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='basic service: task and queue service';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `services_task`
--

LOCK TABLES `services_task` WRITE;
/*!40000 ALTER TABLE `services_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `services_task` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

DROP TABLE IF EXISTS `session`;
CREATE TABLE `session` (
    `session_key` CHAR(64) NOT NULL,
    `session_data` BLOB,
    `session_expiry` int(11) UNSIGNED NOT NULL,
    PRIMARY KEY (`session_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT 'session table, store user sessions';

-- Dump completed on 2019-05-30  8:39:36

CREATE USER 'brunton'@'%' IDENTIFIED BY 'brunton';
ALTER USER 'brunton'@'%' IDENTIFIED WITH mysql_native_password BY 'brunton';
GRANT ALL PRIVILEGES ON b_core.* TO 'brunton'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES ;