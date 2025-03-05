CREATE TABLE `order` (
    -- 订单ID
    `id` BIGINT AUTO_INCREMENT NOT NULL,
    -- 订单hash
    `order_hash` VARCHAR(255) NOT NULL,
    -- 0 等待成交/ 1 已成交/ 2 部分成交/3 已取消/4 已失败
    `status` INT NOT NULL DEFAULT 0,
    -- 订单提示
    `message` VARCHAR(1024),  -- 假设最大长度为1024
    -- 链ID
    `chain_id` VARCHAR(255) NOT NULL,
    -- 市场地址
    `market_address` VARCHAR(255) NOT NULL,
    -- 0buy, 1sell
    `side` INT NOT NULL,
    -- 0市价单/1限价单
    `type` INT NOT NULL DEFAULT 0,
    -- 0市值/1百分比/2价格涨跌幅
    `limit_order_type` INT NOT NULL DEFAULT 0,
    -- 成交时价格
    `price` DOUBLE NOT NULL,
    -- 买入/卖出数量
    `amount` DOUBLE NOT NULL,
    -- 滑点
    `slippage` DOUBLE NOT NULL DEFAULT 0,
    -- 已成交数量
    `filled_amount` DOUBLE NOT NULL DEFAULT 0,
    -- 剩余数量
    `remaining_amount` DOUBLE NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    -- 用户ID
    `account_id` BIGINT NOT NULL,
    -- 0未执行/1正在上链/2已上链/3上链失败
    `payment_status` INT NOT NULL DEFAULT 0,
    -- 交易凭证
    `transaction_hash` VARCHAR(255),  -- 假设最大长度为255
    -- 取消原因
    `cancel_reason` VARCHAR(1024),  -- 假设最大长度为1024
    -- 是否开启防夹模式
    `open_mev` TINYINT(1) NOT NULL,
    -- 0未返佣/1已返佣
    `rebate_status` INT NOT NULL DEFAULT 0,
    PRIMARY KEY(`id`),
    UNIQUE KEY `order_hash_unique` (`order_hash`)
);
