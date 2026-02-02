CREATE TABLE IF NOT EXISTS order_items (
    order_item_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id         UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    sku_id           VARCHAR(64) NOT NULL,
    quantity         INT NOT NULL,
    price            DECIMAL(10,2) NOT NULL,
    currency         CHAR(3) NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_order_items_quantity CHECK (quantity > 0)
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);