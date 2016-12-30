/*
Navicat MySQL Data Transfer

Source Server         : 127.0.0.1
Source Server Version : 50520
Source Host           : localhost:3306
Source Database       : hcloud

Target Server Type    : MYSQL
Target Server Version : 50520
File Encoding         : 65001

Date: 2016-12-19 20:31:54
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for cloud_distribut_record
-- ----------------------------
DROP TABLE IF EXISTS `cloud_distribut_record`;
CREATE TABLE `cloud_distribut_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `down_number` bigint(20) NOT NULL DEFAULT '0',
  `distribut_name` bigint(20) NOT NULL DEFAULT '0',
  `distri_resource_type` bigint(20) NOT NULL DEFAULT '0',
  `down_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cloud_distribut_record
-- ----------------------------

-- ----------------------------
-- Table structure for cloud_resources
-- ----------------------------
DROP TABLE IF EXISTS `cloud_resources`;
CREATE TABLE `cloud_resources` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `phone_number` bigint(20) NOT NULL DEFAULT '0',
  `cloud_resource_type_id` bigint(20) NOT NULL,
  `cloud_upload_record_id` bigint(20) NOT NULL,
  `allocation_status` int(11) NOT NULL DEFAULT '0',
  `downtime` bigint(20) NOT NULL DEFAULT '0',
  `distri_team` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_number` (`phone_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cloud_resources
-- ----------------------------

-- ----------------------------
-- Table structure for cloud_resource_record
-- ----------------------------
DROP TABLE IF EXISTS `cloud_resource_record`;
CREATE TABLE `cloud_resource_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `resource_type` bigint(20) NOT NULL DEFAULT '0',
  `source_total` bigint(20) NOT NULL DEFAULT '0',
  `used_source` bigint(20) NOT NULL DEFAULT '0',
  `canuse_resource` bigint(20) NOT NULL DEFAULT '0',
  `user_persent` varchar(255) NOT NULL DEFAULT '',
  `add_time` date NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cloud_resource_record
-- ----------------------------

-- ----------------------------
-- Table structure for cloud_resource_total
-- ----------------------------
DROP TABLE IF EXISTS `cloud_resource_total`;
CREATE TABLE `cloud_resource_total` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `resource_date` date NOT NULL,
  `total_resource` bigint(20) NOT NULL DEFAULT '0',
  `usde_resource` bigint(20) NOT NULL DEFAULT '0',
  `can_use_resource` bigint(20) NOT NULL DEFAULT '0',
  `use_persent` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cloud_resource_total
-- ----------------------------

-- ----------------------------
-- Table structure for cloud_resource_type
-- ----------------------------
DROP TABLE IF EXISTS `cloud_resource_type`;
CREATE TABLE `cloud_resource_type` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `type_name` varchar(255) NOT NULL DEFAULT '',
  `add_uers` varchar(255) NOT NULL DEFAULT '',
  `add_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `type_name` (`type_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cloud_resource_type
-- ----------------------------

-- ----------------------------
-- Table structure for cloud_upload_record
-- ----------------------------
DROP TABLE IF EXISTS `cloud_upload_record`;
CREATE TABLE `cloud_upload_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `cloud_resource_type_id` bigint(20) NOT NULL,
  `repat_number` bigint(20) NOT NULL DEFAULT '0',
  `success_number` bigint(20) NOT NULL DEFAULT '0',
  `failued_number` bigint(20) NOT NULL DEFAULT '0',
  `total_number` bigint(20) NOT NULL DEFAULT '0',
  `upload_status` int(11) NOT NULL DEFAULT '0',
  `upload_name` varchar(255) NOT NULL DEFAULT '',
  `upload_user` varchar(255) NOT NULL DEFAULT '',
  `upload_date` datetime NOT NULL,
  `bar` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of cloud_upload_record
-- ----------------------------

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `status` bigint(20) NOT NULL DEFAULT '0',
  `fid` bigint(20) NOT NULL DEFAULT '0',
  `sort` bigint(20) NOT NULL DEFAULT '0',
  `remark` varchar(255) NOT NULL DEFAULT '',
  `tid` bigint(20) NOT NULL DEFAULT '0',
  `down_one` bigint(20) NOT NULL DEFAULT '0',
  `one_day_limit` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of group
-- ----------------------------
INSERT INTO `group` VALUES ('1', '管理组', '1', '1', '0', '0', '系统初始默认组', '0', '0', '0');
INSERT INTO `group` VALUES ('2', '中天鑫业', '2', '1', '1', '50', '武汉皓月浮光科技', '0', '0', '0');
INSERT INTO `group` VALUES ('3', '微信运营部', '3', '1', '2', '50', '微信运营部', '2', '0', '0');
INSERT INTO `group` VALUES ('4', '微信运营一部', '5', '1', '3', '50', '微信运营一部', '3', '0', '0');
INSERT INTO `group` VALUES ('5', '微信运营二部', '6', '1', '3', '50', '微信运营二部', '0', '0', '0');
INSERT INTO `group` VALUES ('6', '微信运营三部', '7', '1', '3', '50', '微信运营三部', '3', '0', '0');
INSERT INTO `group` VALUES ('7', '业务拓展部', '4', '1', '2', '50', '业务拓展部', '2', '0', '0');
INSERT INTO `group` VALUES ('8', '业务拓展一部', '0', '1', '7', '50', '业务拓展一部', '7', '0', '0');
INSERT INTO `group` VALUES ('9', '西安恒通嘉诚', '17', '1', '1', '50', '西安恒通嘉诚', '0', '0', '0');
INSERT INTO `group` VALUES ('10', '市场一部', '19', '1', '9', '50', '市场一部', '9', '0', '0');
INSERT INTO `group` VALUES ('11', '市场二部', '20', '1', '9', '50', '市场二部', '9', '0', '0');
INSERT INTO `group` VALUES ('12', '市场三部', '0', '1', '9', '50', '市场三部', '9', '0', '0');
INSERT INTO `group` VALUES ('13', '市场一部一组', '21', '1', '10', '50', '市场一部一组', '10', '0', '0');
INSERT INTO `group` VALUES ('14', '市场一部二组', '22', '1', '10', '50', '市场一部二组', '10', '0', '0');
INSERT INTO `group` VALUES ('15', '市场二部一组', '23', '1', '11', '50', '市场二部一组', '11', '0', '0');
INSERT INTO `group` VALUES ('16', '上海积木', '33', '1', '1', '50', '上海积木', '0', '0', '0');
INSERT INTO `group` VALUES ('17', '积木一组', '0', '1', '16', '50', '积木一组', '16', '0', '0');
INSERT INTO `group` VALUES ('18', '积木二组', '1', '1', '16', '50', '积木二组', '16', '0', '0');

-- ----------------------------
-- Table structure for mgr_device
-- ----------------------------
DROP TABLE IF EXISTS `mgr_device`;
CREATE TABLE `mgr_device` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `uuid` varchar(255) NOT NULL DEFAULT '',
  `serial` varchar(255) NOT NULL DEFAULT '',
  `usb` varchar(255) NOT NULL DEFAULT '',
  `model` varchar(255) NOT NULL DEFAULT '',
  `sdk_version` varchar(255) NOT NULL DEFAULT '',
  `version` varchar(255) NOT NULL DEFAULT '',
  `abi` varchar(255) NOT NULL DEFAULT '',
  `height` int(11) NOT NULL DEFAULT '0',
  `width` int(11) NOT NULL DEFAULT '0',
  `imei` varchar(255) NOT NULL DEFAULT '',
  `online` int(11) NOT NULL DEFAULT '0',
  `nick_name` varchar(255) NOT NULL DEFAULT '',
  `order` int(11) NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL,
  `create_time` datetime NOT NULL,
  `remark` longtext NOT NULL,
  `group_id` bigint(20) NOT NULL DEFAULT '0',
  `muid` bigint(20) NOT NULL DEFAULT '0',
  `syslist_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `syslist_id` (`syslist_id`),
  KEY `mgr_device_usb` (`usb`),
  KEY `mgr_device_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of mgr_device
-- ----------------------------

-- ----------------------------
-- Table structure for mgr_group
-- ----------------------------
DROP TABLE IF EXISTS `mgr_group`;
CREATE TABLE `mgr_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of mgr_group
-- ----------------------------

-- ----------------------------
-- Table structure for mgr_user_device
-- ----------------------------
DROP TABLE IF EXISTS `mgr_user_device`;
CREATE TABLE `mgr_user_device` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `fid` bigint(20) NOT NULL DEFAULT '0',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `did` bigint(20) NOT NULL DEFAULT '0',
  `to_uid` bigint(20) NOT NULL DEFAULT '0',
  `statues` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of mgr_user_device
-- ----------------------------

-- ----------------------------
-- Table structure for node
-- ----------------------------
DROP TABLE IF EXISTS `node`;
CREATE TABLE `node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL DEFAULT '',
  `pid` bigint(20) NOT NULL DEFAULT '0',
  `key` varchar(64) NOT NULL DEFAULT '',
  `type` varchar(10) NOT NULL DEFAULT '',
  `ico` varchar(255) NOT NULL DEFAULT '',
  `url` varchar(64) NOT NULL DEFAULT '',
  `fid` bigint(20) NOT NULL DEFAULT '0',
  `level` bigint(20) NOT NULL DEFAULT '1',
  `description` varchar(200) DEFAULT NULL,
  `sort` bigint(20) NOT NULL DEFAULT '100',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of node
-- ----------------------------
INSERT INTO `node` VALUES ('1', '用户管理', '0', 'user_list', '0', 'fa-user', 'hcloud/user/list', '0', '1', '用户管理', '0');
INSERT INTO `node` VALUES ('2', '职务管理', '0', 'role', '0', 'fa-flag', 'hcloud/role/list', '0', '1', '职务管理', '0');
INSERT INTO `node` VALUES ('3', '团队管理', '0', 'group', '0', 'fa-users', 'hcloud/group/list', '0', '1', '团队管理', '0');
INSERT INTO `node` VALUES ('4', '新增用户', '1000', 'user_add', '2', '', 'hcloud/user/add', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('5', '编辑用户', '1000', 'user_edit', '2', '', 'hcloud/user/edit', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('6', '删除用户', '1000', 'user_del', '2', '', 'hcloud/user/delete', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('7', '用户查看', '1000', 'user_view', '2', '', 'hcloud/user/view', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('8', '用户状态', '1000', 'user_edit_status', '2', '', 'hcloud/user/edituserstatus', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('9', '用户导出', '1000', 'user_list_out', '2', '', 'hcloud/user/export', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('10', '修改密码', '1000', 'user_editpwd', '2', '', 'hcloud/user/editpwd', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('11', '登录记录', '1000', 'user_login_record', '2', '', 'hcloud/user/getLoginRecord', '1', '2', '用户管理-新增用户', '1');
INSERT INTO `node` VALUES ('12', '新增职务', '1000', 'role_add', '2', '', 'hcloud/role/add', '2', '2', '职务管理-新增职务', '2');
INSERT INTO `node` VALUES ('13', '编辑职务', '1000', 'role_edit', '2', '', 'hcloud/role/edit', '2', '2', '职务管理-新增职务', '2');
INSERT INTO `node` VALUES ('14', '删除职务', '1000', 'role_del', '2', '', 'hcloud/role/delete', '2', '2', '职务管理-新增职务', '2');
INSERT INTO `node` VALUES ('15', '分配权限', '1000', 'role_allocat_permission', '2', '', 'hcloud/role/permission', '2', '2', '职务管理-新增职务', '2');
INSERT INTO `node` VALUES ('16', '新增团队', '1000', 'group_add', '2', '', 'hcloud/group/add', '3', '2', '团队管理-新增团队', '3');
INSERT INTO `node` VALUES ('17', '删除团队', '1000', 'group_del', '2', '', 'hcloud/group/delete', '3', '2', '团队管理-新增团队', '3');
INSERT INTO `node` VALUES ('18', '编辑团队', '1000', 'group_edit', '2', '', 'hcloud/group/edit', '3', '2', '团队管理-新增团队', '3');
INSERT INTO `node` VALUES ('19', '设备管理', '0', 'device_list', '0', 'fa-laptop', 'hcloud/device/list', '0', '1', '设备管理', '0');
INSERT INTO `node` VALUES ('20', '分配设备', '1000', 'devie_allot', '2', '', 'hcloud/device/allot', '19', '2', '设备管理-分配设备', '1');


-- ----------------------------
-- Table structure for node_roles
-- ----------------------------
DROP TABLE IF EXISTS `node_roles`;
CREATE TABLE `node_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `node_id` bigint(20) NOT NULL,
  `role_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of node_roles
-- ----------------------------

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `status` bigint(20) NOT NULL DEFAULT '1',
  `description` longtext NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES ('1', '管理员', '1', '系统默认管理员');
INSERT INTO `role` VALUES ('2', '云控账户', '1', '云控账户');
INSERT INTO `role` VALUES ('3', '总经理', '1', '总经理');
INSERT INTO `role` VALUES ('4', '市场总监', '1', '市场总监');
INSERT INTO `role` VALUES ('5', '部门经理', '1', '部门经理');
INSERT INTO `role` VALUES ('6', '小组组长', '1', '小组组长');
INSERT INTO `role` VALUES ('7', '业务员', '1', '业务员');

-- ----------------------------
-- Table structure for sys_list
-- ----------------------------
DROP TABLE IF EXISTS `sys_list`;
CREATE TABLE `sys_list` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `token` varchar(255) NOT NULL DEFAULT '',
  `ip` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `sys_list_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of sys_list
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `nickname` varchar(255) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '1',
  `empnum` varchar(255) NOT NULL DEFAULT '',
  `companyname` varchar(255) NOT NULL DEFAULT '',
  `fid` bigint(20) NOT NULL DEFAULT '0',
  `token` varchar(255) NOT NULL DEFAULT '',
  `remark` longtext,
  `role_id` bigint(20) NOT NULL,
  `group_id` bigint(20) NOT NULL,
  `protected` bigint(20) NOT NULL DEFAULT '0',
  `createtime` datetime NOT NULL,
  `lastlogintime` datetime NOT NULL,
  `one_day_limit` bigint(20) NOT NULL DEFAULT '30',
  `once_limit` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', 'admin', '21232f297a57a5a743894a0e4a801fc3', 'superadmin', '1', '10001', '中天鑫业', '0', '', '', '1', '1', '0', '2016-12-14 09:08:44', '2016-12-19 19:16:19', '0', '0');
INSERT INTO `user` VALUES ('2', 'wh001', '25d55ad283aa400af464c76d713c07ad', '丁碧曼', '1', '001', '中天鑫业', '1', '', '丁碧曼', '1', '2', '0', '2016-12-19 20:01:46', '2016-12-19 20:01:46', '0', '0');
INSERT INTO `user` VALUES ('3', 'wh002', '25d55ad283aa400af464c76d713c07ad', '康谷蓝', '1', '002', '中天鑫业', '2', '', '康谷蓝', '3', '2', '0', '2016-12-19 20:02:40', '2016-12-19 20:02:40', '0', '0');
INSERT INTO `user` VALUES ('4', 'wh003', '25d55ad283aa400af464c76d713c07ad', '藩悠馨', '1', '003', '中天鑫业', '2', '', '藩悠馨', '4', '3', '0', '2016-12-19 20:03:47', '2016-12-19 20:03:47', '0', '0');
INSERT INTO `user` VALUES ('5', 'wh004', '25d55ad283aa400af464c76d713c07ad', '倪香菱', '1', '004', '中天鑫业', '2', '', '倪香菱', '4', '7', '0', '2016-12-19 20:04:33', '2016-12-19 20:04:33', '0', '0');
INSERT INTO `user` VALUES ('6', 'wh005', '25d55ad283aa400af464c76d713c07ad', '符凝云', '1', '005', '中天鑫业', '2', '', '符凝云', '5', '4', '0', '2016-12-19 20:05:18', '2016-12-19 20:05:18', '0', '0');
INSERT INTO `user` VALUES ('7', 'wh006', '25d55ad283aa400af464c76d713c07ad', '沙元蝶', '1', '006', '中天鑫业', '2', '', '沙元蝶', '5', '5', '0', '2016-12-19 20:06:07', '2016-12-19 20:06:07', '0', '0');
INSERT INTO `user` VALUES ('8', 'wh007', '25d55ad283aa400af464c76d713c07ad', '生初夏', '1', '007', '中天鑫业', '2', '', '生初夏', '5', '6', '0', '2016-12-19 20:06:53', '2016-12-19 20:06:53', '0', '0');
INSERT INTO `user` VALUES ('9', 'wh008', '25d55ad283aa400af464c76d713c07ad', '余思真', '1', '008', '中天鑫业', '2', '', '余思真', '7', '4', '0', '2016-12-19 20:07:29', '2016-12-19 20:07:29', '0', '0');
INSERT INTO `user` VALUES ('10', 'wh009', '25d55ad283aa400af464c76d713c07ad', '素妙晴', '1', '009', '中天鑫业', '2', '', '素妙晴', '7', '4', '0', '2016-12-19 20:08:11', '2016-12-19 20:08:11', '0', '0');
INSERT INTO `user` VALUES ('11', 'wh010', '25d55ad283aa400af464c76d713c07ad', '颜香梅', '1', '010', '中天鑫业', '2', '', '颜香梅', '7', '5', '0', '2016-12-19 20:08:55', '2016-12-19 20:08:55', '0', '0');
INSERT INTO `user` VALUES ('12', 'wh011', '25d55ad283aa400af464c76d713c07ad', '野可儿', '1', '011', '中天鑫业', '2', '', '野可儿', '7', '5', '0', '2016-12-19 20:09:32', '2016-12-19 20:09:32', '0', '0');
INSERT INTO `user` VALUES ('13', 'wh012', '25d55ad283aa400af464c76d713c07ad', '畅雪卉', '1', '012', '中天鑫业', '2', '', '畅雪卉', '7', '6', '0', '2016-12-19 20:10:33', '2016-12-19 20:10:33', '0', '0');
INSERT INTO `user` VALUES ('14', 'wh013', '25d55ad283aa400af464c76d713c07ad', '斋璇娟', '1', '013', '中天鑫业', '2', '', '斋璇娟', '7', '5', '0', '2016-12-19 20:11:11', '2016-12-19 20:11:11', '0', '0');
INSERT INTO `user` VALUES ('15', 'wh014', '25d55ad283aa400af464c76d713c07ad', '郭才艺', '1', '014', '中天鑫业', '2', '', '郭才艺', '7', '7', '0', '2016-12-19 20:12:45', '2016-12-19 20:12:45', '0', '0');
INSERT INTO `user` VALUES ('16', 'wh015', '25d55ad283aa400af464c76d713c07ad', '禚思洁', '1', '015', '中天鑫业', '2', '', '稽若彤', '7', '8', '0', '2016-12-19 20:13:33', '2016-12-19 20:13:33', '0', '0');
INSERT INTO `user` VALUES ('17', 'xn1001', '25d55ad283aa400af464c76d713c07ad', '稽若彤', '1', '1001', '西安恒通嘉诚', '1', '', '佛月朗', '1', '9', '0', '2016-12-19 20:07:29', '2016-12-19 20:07:29', '0', '0');
INSERT INTO `user` VALUES ('18', 'xn1002', '25d55ad283aa400af464c76d713c07ad', '佛月朗', '1', '1002', '西安恒通嘉诚', '17', '', '盘思聪', '3', '9', '0', '2016-12-19 20:08:11', '2016-12-19 20:08:11', '0', '0');
INSERT INTO `user` VALUES ('19', 'xn1003', '25d55ad283aa400af464c76d713c07ad', '盘思聪', '1', '1003', '西安恒通嘉诚', '17', '', '司半兰', '4', '10', '0', '2016-12-19 20:08:55', '2016-12-19 20:08:55', '0', '0');
INSERT INTO `user` VALUES ('20', 'xn1004', '25d55ad283aa400af464c76d713c07ad', '司半兰', '1', '1004', '西安恒通嘉诚', '17', '', '盍映寒', '4', '11', '0', '2016-12-19 20:09:32', '2016-12-19 20:09:32', '0', '0');
INSERT INTO `user` VALUES ('21', 'xn1005', '25d55ad283aa400af464c76d713c07ad', '盍映寒', '1', '1005', '西安恒通嘉诚', '17', '', '检梓莹', '5', '13', '0', '2016-12-19 20:10:33', '2016-12-19 20:10:33', '0', '0');
INSERT INTO `user` VALUES ('22', 'xn1006', '25d55ad283aa400af464c76d713c07ad', '检梓莹', '1', '1006', '西安恒通嘉诚', '17', '', '真新之', '5', '14', '0', '2016-12-19 20:11:11', '2016-12-19 20:11:11', '0', '0');
INSERT INTO `user` VALUES ('23', 'xn1007', '25d55ad283aa400af464c76d713c07ad', '真新之', '1', '1007', '西安恒通嘉诚', '17', '', '卞听', '5', '15', '0', '2016-12-19 20:12:45', '2016-12-19 20:12:45', '0', '0');
INSERT INTO `user` VALUES ('24', 'xn1008', '25d55ad283aa400af464c76d713c07ad', '卞听', '1', '1008', '西安恒通嘉诚', '17', '', '禚思洁', '7', '13', '0', '2016-12-19 20:13:33', '2016-12-19 20:13:33', '0', '0');
INSERT INTO `user` VALUES ('25', 'xn1009', '25d55ad283aa400af464c76d713c07ad', '稽若彤', '1', '1009', '西安恒通嘉诚', '17', '', '干才英', '7', '13', '0', '2016-12-19 20:07:29', '2016-12-19 20:07:29', '0', '0');
INSERT INTO `user` VALUES ('26', 'xn1010', '25d55ad283aa400af464c76d713c07ad', '佛月朗', '1', '1010', '西安恒通嘉诚', '17', '', '肇半青', '7', '14', '0', '2016-12-19 20:08:11', '2016-12-19 20:08:11', '0', '0');
INSERT INTO `user` VALUES ('27', 'xn1011', '25d55ad283aa400af464c76d713c07ad', '盘思聪', '1', '1011', '西安恒通嘉诚', '17', '', '籍宇达', '7', '15', '0', '2016-12-19 20:08:55', '2016-12-19 20:08:55', '0', '0');
INSERT INTO `user` VALUES ('28', 'xn1012', '25d55ad283aa400af464c76d713c07ad', '司半兰', '1', '1012', '西安恒通嘉诚', '17', '', '凭颐然', '7', '13', '0', '2016-12-19 20:09:32', '2016-12-19 20:09:32', '0', '0');
INSERT INTO `user` VALUES ('29', 'xn1013', '25d55ad283aa400af464c76d713c07ad', '盍映寒', '1', '1013', '西安恒通嘉诚', '17', '', '泥新儿', '7', '14', '0', '2016-12-19 20:10:33', '2016-12-19 20:10:33', '0', '0');
INSERT INTO `user` VALUES ('30', 'xn1014', '25d55ad283aa400af464c76d713c07ad', '检梓莹', '1', '1014', '西安恒通嘉诚', '17', '', '镇宜楠', '7', '13', '0', '2016-12-19 20:11:11', '2016-12-19 20:11:11', '0', '0');
INSERT INTO `user` VALUES ('31', 'xn1015', '25d55ad283aa400af464c76d713c07ad', '真新之', '1', '1015', '西安恒通嘉诚', '17', '', '盖霏霏', '7', '15', '0', '2016-12-19 20:12:45', '2016-12-19 20:12:45', '0', '0');
INSERT INTO `user` VALUES ('32', 'xn1016', '25d55ad283aa400af464c76d713c07ad', '卞听', '1', '1016', '西安恒通嘉诚', '17', '', '黄冷霜', '7', '10', '0', '2016-12-19 20:13:33', '2016-12-19 20:13:33', '0', '0');
INSERT INTO `user` VALUES ('33', '8001', '25d55ad283aa400af464c76d713c07ad', '8001', '1', '8001', '上海积木', '1', '', '8001', '1', '16', '0', '2016-12-19 20:30:24', '2016-12-19 20:30:24', '0', '0');
INSERT INTO `user` VALUES ('34', '8002', '25d55ad283aa400af464c76d713c07ad', '8002', '1', '8002', '上海积木', '1', '', '上海积木', '7', '17', '0', '2016-12-19 20:30:51', '2016-12-19 20:30:51', '0', '0');
