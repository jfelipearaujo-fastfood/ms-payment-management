DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_items;

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

CREATE TABLE IF NOT EXISTS payment_items (
    id varchar(255),
    order_id varchar(255),
    payment_id varchar(255),
    name varchar(255),
    quantity int,
    PRIMARY KEY (id, order_id, payment_id)
);