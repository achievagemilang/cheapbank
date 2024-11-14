CREATE TYPE Currency AS ENUM (
  'USD',
  'EUR'
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" Currency NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "amount" bigint NOT NULL,
  "account_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "amount" bigint NOT NULL,
  "currency" Currency NOT NULL,
  "from_acc_id" bigint NOT NULL,
  "to_acc_id" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_acc_id");

CREATE INDEX ON "transfers" ("to_acc_id");

CREATE INDEX ON "transfers" ("from_acc_id", "to_acc_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be positive/negative';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_acc_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_acc_id") REFERENCES "accounts" ("id");
