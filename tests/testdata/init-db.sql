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
    PRIMARY KEY (id)
);

INSERT INTO payments(
	order_id, payment_id, total_items, amount, state, created_at, updated_at)
	VALUES (
        'c3fdab1b-3c06-4db2-9edc-4760a2429460',
        '9dfa1386-2f52-4cca-b9aa-f9bd6887d442', 
        1, 
        100.00, 
        1, 
        NOW(), 
        NOW());

INSERT INTO payment_items(
	id, order_id, payment_id, name, quantity)
	VALUES (
        'cfdab175-1f86-4fb0-9bcb-15f2c58df30c',
        'c3fdab1b-3c06-4db2-9edc-4760a2429460',
        '9dfa1386-2f52-4cca-b9aa-f9bd6887d442',
        'Hamburger',
        1);