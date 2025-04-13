CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "customer" (
  "id" UUID UNIQUE PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "phone" VARCHAR NOT NULL,
  "cpf_cnpj" VARCHAR NOT NULL,
  "is_active" BOOLEAN NOT NULL DEFAULT (false),
  "updated_at" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "hashed_password" VARCHAR NOT NULL
);

CREATE INDEX idx_customer_name ON customer (name);
CREATE INDEX idx_customer_email ON customer (email);
CREATE INDEX idx_customer_phone ON customer (phone);
CREATE INDEX idx_customer_updated_at ON customer (updated_at);

CREATE UNIQUE INDEX idx_customer_email_active ON customer (email)
  WHERE is_active = true;

CREATE UNIQUE INDEX idx_customer_phone_active ON customer (phone)
  WHERE is_active = true;