---------------------------------------------------------- DELETE -------------------------------------------------------

-- DROP MATERIALIZED VIEW market_kline_1m;
-- SELECT remove_continuous_aggregate_policy('market_kline_1m');

---------------------------------------------------------- CREATE -------------------------------------------------------


CREATE MATERIALIZED VIEW market_kline_1m
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 minute', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;
SELECT add_continuous_aggregate_policy('market_kline_1m',
                                       start_offset => INTERVAL '3 minute',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '1 minute');







CREATE MATERIALIZED VIEW market_kline_5m
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('5 minute', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;
SELECT add_continuous_aggregate_policy('market_kline_5m',
                                       start_offset => INTERVAL '15 minute',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '5 minute');







CREATE MATERIALIZED VIEW market_kline_15m
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('15 minute', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;
SELECT add_continuous_aggregate_policy('market_kline_15m',
                                       start_offset => INTERVAL '45 minute',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '15 minute');







CREATE MATERIALIZED VIEW market_kline_30m
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('30 minute', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;
SELECT add_continuous_aggregate_policy('market_kline_30m',
                                       start_offset => INTERVAL '90 minute',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '30 minute');




CREATE MATERIALIZED VIEW market_kline_1h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;
SELECT add_continuous_aggregate_policy('market_kline_1h',
                                       start_offset => INTERVAL '3 hour',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '1 hour');








CREATE MATERIALIZED VIEW market_kline_4h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('4 hour', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;
SELECT add_continuous_aggregate_policy('market_kline_4h',
                                       start_offset => INTERVAL '12 hour',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '4 hour');







CREATE MATERIALIZED VIEW market_kline_6h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('6 hour', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;

SELECT add_continuous_aggregate_policy('market_kline_6h',
                                       start_offset => INTERVAL '18 hour',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '6 hour');




CREATE MATERIALIZED VIEW market_kline_12h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('12 hour', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;

SELECT add_continuous_aggregate_policy('market_kline_12h',
                                       start_offset => INTERVAL '36 hour',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '12 hour');



CREATE MATERIALIZED VIEW market_kline_1d
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('24 hour', created_timestamp) AS timestamp,
        market_address,
        FIRST(quote_price, created_timestamp) AS "o",
        MAX(quote_price) AS h,
        MIN(quote_price) AS l,
        LAST(quote_price, created_timestamp) AS "c",
        SUM(volume) AS v
FROM market_tx
GROUP BY timestamp, market_address;

SELECT add_continuous_aggregate_policy('market_kline_1d',
                                       start_offset => INTERVAL '48 hour',
                                       end_offset => NULL,
                                       schedule_interval => INTERVAL '24 hour');


