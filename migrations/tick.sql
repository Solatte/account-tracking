CREATE TABLE IF NOT EXISTS tick (
    amm_id VARCHAR(255),
    open_price BIGINT,
    close_price BIGINT,
    high_price BIGINT,
    low_price BIGINT,
    change_percent FLOAT,
    price_difference BIGINT,
    trend INT,
    timestamp INT
);