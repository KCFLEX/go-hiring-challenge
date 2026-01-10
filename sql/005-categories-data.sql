-- Seed categories
INSERT INTO categories (code, name) VALUES
  ('CAT001', 'Clothing'),
  ('CAT002', 'Shoes'),
  ('CAT003', 'Accessories');

-- Backfill product category_id
UPDATE products SET category_id = c.id
FROM categories c
WHERE
  (products.code IN ('PROD001','PROD004','PROD007') AND c.code = 'CAT001') OR
  (products.code IN ('PROD002','PROD006') AND c.code = 'CAT002') OR
  (products.code IN ('PROD003','PROD005','PROD008') AND c.code = 'CAT003');

-- Enforce NOT NULL after backfill
ALTER TABLE products ALTER COLUMN category_id SET NOT NULL;

-- Optional index for filtering
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products (category_id);