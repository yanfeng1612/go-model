#交易账号信息表  (分表key : s_code):
#   表说明:交易账号由基金云系统产生,记录交易账号基本信息,增卡时写入数据
# 逻辑主键:交易账号
# @{"tableName":"acc_trans_account","shardingKey":["s_code"],"singleQueryField":["transaction_account_id"]}

CREATE TABLE acc_trans_account(
   id BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '自动编号',
   s_code VARCHAR(16) NOT NULL  COMMENT '切分键',
   uuid VARCHAR(100) NOT NULL COMMENT '防重标识 custNo_branchCode_depositAcct',
   trader_no VARCHAR(8) NOT NULL  COMMENT '交易员编号',
   transaction_account_id VARCHAR(100) COMMENT '交易账号',
   cust_no VARCHAR(20) COMMENT '客户号',
   channel_id  VARCHAR(100) COMMENT '支付渠道号',
   branch_code VARCHAR(50) COMMENT '支付网点号',
   deposit_acct VARCHAR(64) NULL COMMENT '系统外资金账号(如果是券商则为资金账号,如果是银行则为银行账号或卡号)',
   yn SMALLINT(2) NOT NULL COMMENT '数据状态 0-无效 1-有效',
   create_pin VARCHAR(100) NOT NULL COMMENT '创建人的pin',
   update_pin VARCHAR(100) NOT NULL COMMENT '修改人的pin',
   created_time DATETIME NOT NULL COMMENT '创建时间',
   modified_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
   PRIMARY KEY(id),
   UNIQUE KEY uniq_uuid (uuid),
   UNIQUE KEY uniq_accountId(transaction_account_id),
   KEY idx_s_code(s_code)
)ENGINE=InnoDB AUTO_INCREMENT= 1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='交易账号信息表';