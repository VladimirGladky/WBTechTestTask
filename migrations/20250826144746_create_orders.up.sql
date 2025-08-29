CREATE TABLE IF NOT EXISTS orders (
                        order_uid VARCHAR(50) PRIMARY KEY,
                        track_number VARCHAR(50),
                        entry VARCHAR(10),
                        locale VARCHAR(10),
                        internal_signature VARCHAR(50),
                        customer_id VARCHAR(50),
                        delivery_service VARCHAR(50),
                        shardkey VARCHAR(10),
                        sm_id INTEGER,
                        date_created VARCHAR(50),
                        oof_shard VARCHAR(10),
                        delivery JSONB,
                        payment JSONB,
                        items JSONB
);

CREATE INDEX idx_orders_order_uid ON orders(order_uid);