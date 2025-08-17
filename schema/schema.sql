-- Создание таблицы deliveries
CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL
);

-- Создание таблицы payments
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(100) NOT NULL UNIQUE,
    request_id VARCHAR(100),
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT DEFAULT 0
);

-- Создание таблицы items
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    chrt_id INT NOT NULL,
    track_number VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INT NOT NULL DEFAULT 0,
    size VARCHAR(20) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INT NOT NULL
);

-- Создание таблицы orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(100) NOT NULL UNIQUE,
    track_number VARCHAR(100) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    delivery_id INT REFERENCES deliveries(id) ON DELETE CASCADE,
    payment_id INT REFERENCES payments(id) ON DELETE CASCADE,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(100),
    customer_id VARCHAR(100) NOT NULL,
    delivery_service VARCHAR(100) NOT NULL,
    shard_key VARCHAR(10) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard VARCHAR(10) NOT NULL
);

-- Создание таблицы связи order_items
CREATE TABLE order_items (
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    item_id INT REFERENCES items(id) ON DELETE CASCADE,
    PRIMARY KEY (order_id, item_id)
);