-- Drop publication
DROP PUBLICATION IF EXISTS payment_dbz_publication;

-- Drop indexes (indexes are automatically dropped with the table, 
-- but we explicitly drop in case of specific use)
DROP INDEX IF EXISTS payments_order_id_idx;
DROP INDEX IF EXISTS payments_transaction_id_idx;
DROP INDEX IF EXISTS outboxes_aggregated_type_aggregated_id_idx;

-- Drop tables
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS outboxes;

-- Drop enums
DROP TYPE IF EXISTS payment_statuses;
DROP TYPE IF EXISTS payment_methods;
