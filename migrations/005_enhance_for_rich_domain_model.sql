-- +goose Up

-- Enhance orders table for rich domain model
-- +goose StatementBegin
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS merchant_id UUID,
ADD COLUMN IF NOT EXISTS total_amount_satoshis BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS delivery_method VARCHAR(50),
ADD COLUMN IF NOT EXISTS delivery_address JSONB,
ADD COLUMN IF NOT EXISTS estimated_window JSONB,
ADD COLUMN IF NOT EXISTS status_history JSONB DEFAULT '[]'::jsonb,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- Update existing orders table to use string status instead of int
-- +goose StatementBegin
-- First add a temporary column for the new status
ALTER TABLE orders ADD COLUMN IF NOT EXISTS status_new VARCHAR(50);

-- Update the new column based on old integer values
UPDATE orders SET status_new = CASE 
    WHEN status = 0 THEN 'PENDING'
    WHEN status = 1 THEN 'ACCEPTED'  
    WHEN status = 2 THEN 'PREPARING'
    WHEN status = 3 THEN 'COMPLETED'
    ELSE 'PENDING'
END WHERE status_new IS NULL;

-- Drop the old column and rename the new one
ALTER TABLE orders DROP COLUMN IF EXISTS status;
ALTER TABLE orders RENAME COLUMN status_new TO status;

-- Set NOT NULL constraint
ALTER TABLE orders ALTER COLUMN status SET NOT NULL;
-- +goose StatementEnd

-- Enhance order_items table for rich domain model
-- +goose StatementBegin
ALTER TABLE order_items
ADD COLUMN IF NOT EXISTS menu_item_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS unit_price_satoshis BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS subtotal_price_satoshis BIGINT DEFAULT 0;

-- Update order_items id to UUID for consistency
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS id_new UUID DEFAULT gen_random_uuid();
UPDATE order_items SET id_new = gen_random_uuid() WHERE id_new IS NULL;
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS order_items_pkey;
ALTER TABLE order_items DROP COLUMN IF EXISTS id;
ALTER TABLE order_items RENAME COLUMN id_new TO id;
ALTER TABLE order_items ADD PRIMARY KEY (id);
-- +goose StatementEnd

-- Enhance menu_items table for rich domain model  
-- +goose StatementBegin
ALTER TABLE menu_items
ADD COLUMN IF NOT EXISTS category VARCHAR(100) DEFAULT 'General',
ADD COLUMN IF NOT EXISTS is_available BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS stock INTEGER DEFAULT -1, -- -1 means unlimited
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Rename price to price_satoshis for clarity
ALTER TABLE menu_items RENAME COLUMN price TO price_satoshis;

-- Drop the old in_stock column and use is_available instead
ALTER TABLE menu_items DROP COLUMN IF EXISTS in_stock;

-- Ensure NOT NULL constraints
ALTER TABLE menu_items ALTER COLUMN is_available SET NOT NULL;
ALTER TABLE menu_items ALTER COLUMN stock SET NOT NULL;
-- +goose StatementEnd

-- Enhance merchants table for rich domain model
-- +goose StatementBegin
ALTER TABLE merchants
ADD COLUMN IF NOT EXISTS operating_hours JSONB,
ADD COLUMN IF NOT EXISTS preparation_time INTEGER DEFAULT 30,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Ensure NOT NULL constraints
ALTER TABLE merchants ALTER COLUMN preparation_time SET NOT NULL;
-- +goose StatementEnd

-- Add foreign key constraints for better data integrity
-- +goose StatementBegin
-- Add foreign key from orders to merchants (if merchants table exists)
ALTER TABLE orders 
ADD CONSTRAINT IF NOT EXISTS fk_orders_merchant_id 
FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE CASCADE;

-- Add foreign key from menu_items to merchants  
ALTER TABLE menu_items
ADD CONSTRAINT IF NOT EXISTS fk_menu_items_merchant_id
FOREIGN KEY (merchant_id) REFERENCES merchants(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- Create indexes for efficient querying
-- +goose StatementBegin
-- Orders indexes
CREATE INDEX IF NOT EXISTS idx_orders_merchant_id ON orders(merchant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_merchant_status ON orders(merchant_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);

-- Order items indexes  
CREATE INDEX IF NOT EXISTS idx_order_items_menu_item_id ON order_items(menu_item_id);

-- Menu items indexes
CREATE INDEX IF NOT EXISTS idx_menu_items_available ON menu_items(is_available);
CREATE INDEX IF NOT EXISTS idx_menu_items_merchant_available ON menu_items(merchant_id, is_available);
CREATE INDEX IF NOT EXISTS idx_menu_items_category ON menu_items(category);

-- Merchants indexes
CREATE INDEX IF NOT EXISTS idx_merchants_active ON merchants(is_active);
-- +goose StatementEnd

-- Create trigger to update updated_at columns automatically
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to tables with updated_at columns
CREATE TRIGGER IF NOT EXISTS update_orders_updated_at 
    BEFORE UPDATE ON orders 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER IF NOT EXISTS update_menu_items_updated_at 
    BEFORE UPDATE ON menu_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER IF NOT EXISTS update_merchants_updated_at 
    BEFORE UPDATE ON merchants 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- Add check constraints for data validation
-- +goose StatementBegin
-- Order status must be valid
ALTER TABLE orders 
ADD CONSTRAINT IF NOT EXISTS chk_orders_status 
CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED', 'PREPARING', 'READY', 'OUT_FOR_DELIVERY', 'COMPLETED', 'CANCELLED'));

-- Delivery method must be valid
ALTER TABLE orders
ADD CONSTRAINT IF NOT EXISTS chk_orders_delivery_method
CHECK (delivery_method IS NULL OR delivery_method IN ('PICKUP', 'DELIVERY'));

-- Order amounts must be positive
ALTER TABLE orders
ADD CONSTRAINT IF NOT EXISTS chk_orders_total_amount_positive
CHECK (total_amount_satoshis >= 0);

-- Order item quantities must be positive
ALTER TABLE order_items
ADD CONSTRAINT IF NOT EXISTS chk_order_items_quantity_positive  
CHECK (quantity > 0);

-- Order item prices must be non-negative
ALTER TABLE order_items
ADD CONSTRAINT IF NOT EXISTS chk_order_items_unit_price_positive
CHECK (unit_price_satoshis >= 0);

ALTER TABLE order_items  
ADD CONSTRAINT IF NOT EXISTS chk_order_items_subtotal_positive
CHECK (subtotal_price_satoshis >= 0);

-- Menu item prices must be positive
ALTER TABLE menu_items
ADD CONSTRAINT IF NOT EXISTS chk_menu_items_price_positive
CHECK (price_satoshis > 0);

-- Stock must be -1 (unlimited) or positive
ALTER TABLE menu_items
ADD CONSTRAINT IF NOT EXISTS chk_menu_items_stock_valid
CHECK (stock = -1 OR stock >= 0);

-- Merchant preparation time must be positive
ALTER TABLE merchants
ADD CONSTRAINT IF NOT EXISTS chk_merchants_prep_time_positive
CHECK (preparation_time > 0);
-- +goose StatementEnd

-- +goose Down

-- Remove check constraints
-- +goose StatementBegin
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_status;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_delivery_method;  
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_total_amount_positive;

ALTER TABLE order_items DROP CONSTRAINT IF EXISTS chk_order_items_quantity_positive;
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS chk_order_items_unit_price_positive;
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS chk_order_items_subtotal_positive;

ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS chk_menu_items_price_positive;
ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS chk_menu_items_stock_valid;

ALTER TABLE merchants DROP CONSTRAINT IF EXISTS chk_merchants_prep_time_positive;
-- +goose StatementEnd

-- Remove triggers and functions
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS update_menu_items_updated_at ON menu_items;
DROP TRIGGER IF EXISTS update_merchants_updated_at ON merchants;
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd

-- Remove indexes
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_merchant_id;
DROP INDEX IF EXISTS idx_orders_status; 
DROP INDEX IF EXISTS idx_orders_merchant_status;
DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_order_items_menu_item_id;
DROP INDEX IF EXISTS idx_menu_items_available;
DROP INDEX IF EXISTS idx_menu_items_merchant_available;
DROP INDEX IF EXISTS idx_menu_items_category;
DROP INDEX IF EXISTS idx_merchants_active;
-- +goose StatementEnd

-- Remove foreign key constraints
-- +goose StatementBegin
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_merchant_id;
ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS fk_menu_items_merchant_id;
-- +goose StatementEnd

-- Revert merchants table changes
-- +goose StatementBegin
ALTER TABLE merchants DROP COLUMN IF EXISTS operating_hours;
ALTER TABLE merchants DROP COLUMN IF EXISTS preparation_time;
ALTER TABLE merchants DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd

-- Revert menu_items table changes
-- +goose StatementBegin
ALTER TABLE menu_items DROP COLUMN IF EXISTS category;
ALTER TABLE menu_items DROP COLUMN IF EXISTS is_available;
ALTER TABLE menu_items DROP COLUMN IF EXISTS stock;
ALTER TABLE menu_items DROP COLUMN IF EXISTS updated_at;
ALTER TABLE menu_items RENAME COLUMN price_satoshis TO price;
ALTER TABLE menu_items ADD COLUMN IF NOT EXISTS in_stock BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd

-- Revert order_items table changes
-- +goose StatementBegin
ALTER TABLE order_items DROP COLUMN IF EXISTS menu_item_name;
ALTER TABLE order_items DROP COLUMN IF EXISTS unit_price_satoshis;
ALTER TABLE order_items DROP COLUMN IF EXISTS subtotal_price_satoshis;
-- Note: Reverting UUID to SERIAL is complex and data-lossy, so we keep UUID but add SERIAL
-- +goose StatementEnd

-- Revert orders table changes
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN IF EXISTS merchant_id;
ALTER TABLE orders DROP COLUMN IF EXISTS total_amount_satoshis;
ALTER TABLE orders DROP COLUMN IF EXISTS delivery_method;
ALTER TABLE orders DROP COLUMN IF EXISTS delivery_address;
ALTER TABLE orders DROP COLUMN IF EXISTS estimated_window;
ALTER TABLE orders DROP COLUMN IF EXISTS status_history;
ALTER TABLE orders DROP COLUMN IF EXISTS updated_at;

-- Revert status back to integer (this will lose data precision)
ALTER TABLE orders ADD COLUMN IF NOT EXISTS status_int INT;
UPDATE orders SET status_int = CASE 
    WHEN status = 'PENDING' THEN 0
    WHEN status = 'ACCEPTED' THEN 1  
    WHEN status = 'PREPARING' THEN 2
    WHEN status = 'COMPLETED' THEN 3
    ELSE 0
END;
ALTER TABLE orders DROP COLUMN IF EXISTS status;
ALTER TABLE orders RENAME COLUMN status_int TO status;
ALTER TABLE orders ALTER COLUMN status SET NOT NULL;
-- +goose StatementEnd
