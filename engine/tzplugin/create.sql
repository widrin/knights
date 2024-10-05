CREATE TABLE `tzp_gift` (
	`id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '礼包id',
	`channel_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '渠道id id1,id2,id3' COLLATE 'utf8mb4_general_ci',
	`name` VARCHAR(255) NULL DEFAULT NULL COMMENT '礼包名字' COLLATE 'utf8mb4_general_ci',
	`condi` INT(11) NULL DEFAULT NULL COMMENT '条件 0:不限制领取次数，1:每天领取一次，2:每周领取一次，3:每月领取一次，4永久领取一次',
	`start_time` DATETIME NULL DEFAULT NULL COMMENT '开始时间',
	`end_time` DATETIME NULL DEFAULT NULL COMMENT '结束时间',
	`title` VARCHAR(255) NULL DEFAULT NULL COMMENT '标题' COLLATE 'utf8mb4_general_ci',
	`content` VARCHAR(1024) NULL DEFAULT NULL COMMENT '内容' COLLATE 'utf8mb4_general_ci',
	`status` INT(11) NULL DEFAULT NULL COMMENT '状态：0禁用，1启用',
	`attach` VARCHAR(1024) NULL DEFAULT NULL COMMENT '道具列表 道具id:道具数量_道具id:道具数量...' COLLATE 'utf8mb4_general_ci',
	`creator` VARCHAR(255) NULL DEFAULT NULL COMMENT '作者' COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=10000
;

CREATE TABLE `tzp_gift_log` (
	`id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '日志id',
	`gift_id` INT(11) NULL DEFAULT NULL COMMENT '礼包id',
	`channel_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '渠道id' COLLATE 'utf8mb4_general_ci',
	`name` VARCHAR(255) NULL DEFAULT NULL COMMENT '礼包名字' COLLATE 'utf8mb4_general_ci',
	`plat_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '平台id' COLLATE 'utf8mb4_general_ci',
	`account_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '账号id' COLLATE 'utf8mb4_general_ci',
	`role_id` INT(11) NULL DEFAULT NULL COMMENT '角色id',
	`send_time` DATETIME NULL DEFAULT NULL COMMENT '发送时间',
	`status` INT(11) NULL DEFAULT NULL COMMENT '状态 0发放失败 1发放成功',
	`creator` VARCHAR(255) NULL DEFAULT NULL COMMENT '作者' COLLATE 'utf8mb4_general_ci',
	`reason` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1000000
;

CREATE TABLE `tzp_gift_record` (
	`dkey` VARCHAR(255) NOT NULL COMMENT 'userid_giftid_时间(d??,w??,m??)' COLLATE 'utf8mb4_general_ci',
	`ts` DATETIME NULL DEFAULT NULL,
	PRIMARY KEY (`dkey`) USING BTREE
)
COMMENT='记录userid_giftid_时间(d??,w??,m??)'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `tzp_wjx` (
	`wjid` INT(11) NOT NULL AUTO_INCREMENT COMMENT '问卷id',
	`channel_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '渠道名称 id1,id2,id3' COLLATE 'utf8mb4_general_ci',
	`reward` INT(11) NULL DEFAULT NULL COMMENT '奖励 tzp_wjx的id',
	`url` VARCHAR(255) NULL DEFAULT NULL COMMENT '问卷url' COLLATE 'utf8mb4_general_ci',
	`content` VARCHAR(255) NULL DEFAULT NULL COMMENT '问卷说明内容' COLLATE 'utf8mb4_general_ci',
	`name` VARCHAR(255) NULL DEFAULT NULL COMMENT '问卷名称' COLLATE 'utf8mb4_general_ci',
	`min_lev` INT(11) NULL DEFAULT NULL COMMENT '最低等级',
	`max_lev` INT(11) NULL DEFAULT NULL COMMENT '最高等级',
	`min_pay` INT(11) NULL DEFAULT NULL COMMENT '最小充值金额 分',
	`max_pay` INT(11) NULL DEFAULT NULL COMMENT '最大充值金额 分',
	`min_open_day` INT(11) NULL DEFAULT NULL COMMENT '最小开服天数',
	`max_open_day` INT(11) NULL DEFAULT NULL COMMENT '最大开服天数',
	`start_time` DATETIME NULL DEFAULT NULL COMMENT '开始时间',
	`end_time` DATETIME NULL DEFAULT NULL COMMENT '结束时间',
	`status` INT(11) NULL DEFAULT NULL COMMENT '状态：0禁用，1启用',
	`creator` VARCHAR(255) NULL DEFAULT NULL COMMENT '作者' COLLATE 'utf8mb4_general_ci',
	`version` BIGINT(20) NULL DEFAULT NULL COMMENT '毫秒时间戳作为版本号',
	PRIMARY KEY (`wjid`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=10000
;
