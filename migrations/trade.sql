CREATE TABLE IF NOT EXISTS trade (
    amm_id VARCHAR(255),
    mint VARCHAR(255),
    action VARCHAR(255),
    compute_limit INT,
    compute_price INT,
    amount BIGINT,
    signature VARCHAR(255),
    timestamp INT,
    tip VARCHAR(255),
    tip_amount INT,
    status VARCHAR(255),
    signer VARCHAR(255)
);