CREATE TABLE IF NOT EXISTS payments (
    order_id varchar(255),
    payment_id varchar(255),
    total_items int,
    amount DECIMAL(10, 2),
    state int,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (order_id, payment_id)
);