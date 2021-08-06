CREATE DATABASE demo;

USE demo;

DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
   id          BIGINT(20)     NOT NULL       AUTO_INCREMENT          COMMENT '自动编号',
   uuid        VARCHAR(128)   NOT NULL                               COMMENT '防重标识uuid',
   user_id     BIGINT(20)     NOT NULL                               COMMENT '用户id',
   third_no    VARCHAR(128)                                          COMMENT '三方订单号',
   state       SMALLINT(4)    NOT NULL                               COMMENT '状态',
   amount      DECIMAL(18,2)  NOT NULL                               COMMENT '金额',
   remark      VARCHAR(512)                                          COMMENT '备注',
   created_time TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP     COMMENT '创建时间',
   modified_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
   PRIMARY KEY(id),
   UNIQUE INDEX uniq_uuid (uuid),
   INDEX  idx_userId(user_id)
) ENGINE=InnoDB AUTO_INCREMENT= 1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='订单表';
