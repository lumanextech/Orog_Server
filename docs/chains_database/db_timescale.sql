show max_connections;
CREATE EXTENSION IF NOT EXISTS timescaledb;




CREATE TABLE "market_tx" (
                             "market_id" bigint,
                             "market_address" text,
                             "quote_amount" double precision,
                             "volume" double precision,
                             "quote_price" double precision ,
                             "tx_type" int,
                             "base_price" double precision,
                             "tx_hash" text,
                             "quote_address" text,
                             "maker_address" text,
                             "base_amount" double precision,
                             "created_timestamp" timestamptz
);

SELECT create_hypertable('market_tx', by_range('created_timestamp'));

