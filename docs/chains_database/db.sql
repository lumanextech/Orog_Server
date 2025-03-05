
CREATE TABLE "market" (
                          "id"  BIGSERIAL PRIMARY KEY,
                          "address" text,
                          "name" text,
                          "symbol" text,
                          "created_timestamp" timestamptz,
                          "open_timestamp" timestamptz, --开盘时间
                          "market_type" text,
                          "base_symbol" text,
                          "base_icon" text,
                          "base_token_address" text,
                          "base_token_mint_address" text,
                          "base_token_decimals" int,
                          "base_token_is_pc_or_token0" boolean,
                          "quote_symbol" text,
                          "quote_icon" text,
                          "quote_token_address" text,
                          "quote_token_mint_address" text,
                          "quote_token_decimals" int,
                          "quote_token_supply" double precision, -- 供应量
                          "maker_address" text,
                          "init_base_vault" double precision,
                          "init_quote_vault" double precision,

                          "quote_vault" double precision, -- 实时代币数量
                          "base_vault" double precision, -- 实时基础代币数量
                          "quote_price" double precision, -- 实时代币价格
                          "base_price" double precision, -- 实时基础代币价格
                          "liquidity" double precision, -- 流动池价值
                          "market_cap" double precision, -- 市值(quote_token_supply*quote_vault)
                          "swaps" int8, -- 总交易笔数
                          "sells" int8, -- 卖出
                          "buys" int8, -- 买入
                          "volume" double precision, --总成交量
                          "holder_count"  int8, -- 持有人数量

                          "price_change_1m" double precision, -- 1m价格变化
                          "price_change_5m" double precision, -- 5m价格变化
                          "price_change_1h" double precision, -- 1h价格变化
                          "price_change_6h" double precision, -- 6h价格变化
                          "price_change_24h" double precision, -- 24h价格变化

                          "sell_volume_1m" double precision, -- 1m卖出成交额
                          "sell_volume_5m" double precision, -- 5m卖出成交额
                          "sell_volume_1h" double precision, -- 1h卖出成交额
                          "sell_volume_6h" double precision, -- 6h卖出成交额
                          "sell_volume_24h" double precision, -- 24h卖出成交额

                          "buy_volume_1m" double precision, -- 1m买入成交额
                          "buy_volume_5m" double precision, -- 5m买入成交额
                          "buy_volume_1h" double precision, -- 1h买入成交额
                          "buy_volume_6h" double precision, -- 6h买入成交额
                          "buy_volume_24h" double precision, -- 24h买入成交额

                          "buy_count_1m" double precision, -- 1m买入笔数
                          "buy_count_5m" double precision, -- 5m买入笔数
                          "buy_count_1h" double precision, -- 1h买入笔数
                          "buy_count_6h" double precision, -- 6h买入笔数
                          "buy_count_24h" double precision, -- 24h买入笔数

                          "sell_count_1m" double precision, -- 1m卖出笔数
                          "sell_count_5m" double precision, -- 5m卖出笔数
                          "sell_count_1h" double precision, -- 1h卖出笔数
                          "sell_count_6h" double precision, -- 6h卖出笔数
                          "sell_count_24h" double precision, -- 24h卖出笔数

                          "updated_at" timestamptz,
                          "created_at" timestamptz
);

CREATE UNIQUE INDEX idx_market_address ON market (address);


CREATE TABLE "market_token" (
                                "id"  BIGSERIAL PRIMARY KEY,
                                "address" text,
                                "symbol" text,
                                "name" text,
                                "icon" text,
                                "total_supply" double precision,
                                "created_timestamp" timestamptz,
                                "updated_at" timestamptz,
                                "created_at" timestamptz,
                                "creator_address" text,
                                "ca_open_source" boolean,
                                "ca_honey_pot" boolean,
                                "ca_renounced" boolean,
                                "ca_can_sell" boolean,
                                "buy_tax" double precision,
                                "sell_tax" double precision,
                                "average_tax" double precision,
                                "highest_tax" double precision
);

CREATE UNIQUE INDEX idx_market_token_address ON market_token (address);


CREATE TABLE "market_tx_uid" (
                                 uid text,
                                 base_amount double precision,
                                 quote_amount double precision
);

CREATE UNIQUE INDEX idx_uid ON market_tx_uid (uid);


CREATE TABLE "market_holder" (
                                 "id"  BIGSERIAL PRIMARY KEY,
                                 "market_id" bigint,
                                 "market_address" text,                                 
                                 "address" text,
                                 "created_timestamp" timestamptz,
                                 "quote_amount" double precision,
                                 "base_amount" double precision,
                                 "created_at" timestamptz
);

CREATE UNIQUE INDEX idx_market_holder_market_address ON market_holder (market_address);
