CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "company" (
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

CREATE INDEX idx_company_name ON company (name);
CREATE INDEX idx_company_email ON company (email);
CREATE INDEX idx_company_phone ON company (phone);
CREATE INDEX idx_company_updated_at ON company (updated_at);

CREATE UNIQUE INDEX idx_company_email_active ON company (email)
  WHERE is_active = true;

CREATE UNIQUE INDEX idx_company_phone_active ON company (phone)
  WHERE is_active = true;
