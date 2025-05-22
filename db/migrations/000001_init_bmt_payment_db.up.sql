CREATE TYPE "payment_statuses" AS ENUM (
    'created',
    'success',
    'failed',
    'canceled',
    'expired'
);

CREATE TYPE "payment_methods" AS ENUM (
    'momo',
    'vnpay',
    'zalopay',
    'credit_card',
    'bank_transfer'
);

CREATE TABLE
    "payments" (
        "id" serial PRIMARY KEY NOT NULL,
        "order_id" int NOT NULL,
        "amount" varchar(32) NOT NULL DEFAULT '0 VND',
        "status" payment_statuses NOT NULL DEFAULT 'created',
        "method" payment_methods NOT NULL,
        "transaction_id" varchar(128) NOT NULL,
        "error_message" text DEFAULT '',
        "created_at" timestamp NOT NULL DEFAULT (now ())
    );

CREATE TABLE
    "outboxes" (
        "id" uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid ()),
        "aggregated_type" varchar(64) NOT NULL,
        "aggregated_id" int NOT NULL,
        "event_type" varchar(64) NOT NULL,
        "payload" jsonb NOT NULL,
        "created_at" timestamp NOT NULL DEFAULT (now ())
    );

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("transaction_id");

CREATE INDEX ON "outboxes" ("aggregated_type", "aggregated_id");

CREATE PUBLICATION payment_dbz_publication FOR TABLE outboxes;