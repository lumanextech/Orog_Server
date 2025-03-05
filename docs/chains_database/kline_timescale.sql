
CREATE TABLE "market_kline_1s" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_1s', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_1s_address_time
    ON market_kline_1s("market_address", "timestamp");


CREATE TABLE "market_kline_1m" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_1m', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_1m_address_time
    ON market_kline_1m("market_address", "timestamp");



CREATE TABLE "market_kline_5m" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_5m', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_5m_address_time
    ON market_kline_5m("market_address", "timestamp");



CREATE TABLE "market_kline_15m" (
                                    "market_address" text,
                                    "o" double precision,
                                    "h" double precision,
                                    "l" double precision,
                                    "c" double precision,
                                    "v" double precision,
                                    "timestamp" timestamptz,
                                    "updated_at" timestamptz,
                                    "created_at" timestamptz
);

SELECT create_hypertable('market_kline_15m', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_15m_address_time
    ON market_kline_15m("market_address", "timestamp");



CREATE TABLE "market_kline_30m" (
                                    "market_address" text,
                                    "o" double precision,
                                    "h" double precision,
                                    "l" double precision,
                                    "c" double precision,
                                    "v" double precision,
                                    "timestamp" timestamptz,
                                    "updated_at" timestamptz,
                                    "created_at" timestamptz
);

SELECT create_hypertable('market_kline_30m', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_30m_address_time
    ON market_kline_30m("market_address", "timestamp");





CREATE TABLE "market_kline_1h" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_1h', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_1h_address_time
    ON market_kline_1h("market_address", "timestamp");








CREATE TABLE "market_kline_4h" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_4h', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_4h_address_time
    ON market_kline_4h("market_address", "timestamp");





CREATE TABLE "market_kline_6h" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_6h', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_6h_address_time
    ON market_kline_6h("market_address", "timestamp");




CREATE TABLE "market_kline_12h" (
                                    "market_address" text,
                                    "o" double precision,
                                    "h" double precision,
                                    "l" double precision,
                                    "c" double precision,
                                    "v" double precision,
                                    "timestamp" timestamptz,
                                    "updated_at" timestamptz,
                                    "created_at" timestamptz
);

SELECT create_hypertable('market_kline_12h', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_12h_address_time
    ON market_kline_12h("market_address", "timestamp");




CREATE TABLE "market_kline_1d" (
                                   "market_address" text,
                                   "o" double precision,
                                   "h" double precision,
                                   "l" double precision,
                                   "c" double precision,
                                   "v" double precision,
                                   "timestamp" timestamptz,
                                   "updated_at" timestamptz,
                                   "created_at" timestamptz
);

SELECT create_hypertable('market_kline_1d', by_range('timestamp'));

CREATE UNIQUE INDEX idx_market_kline_1d_address_time
    ON market_kline_1d("market_address", "timestamp");