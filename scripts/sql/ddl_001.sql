CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE schema if not EXISTS "user";
CREATE schema if not EXISTS loan;
CREATE schema if not EXISTS "file";

CREATE  TYPE loan.loan_state AS ENUM ('proposed', 'approved', 'invested', 'disbursed', 'rejected');
CREATE  TYPE "user".user_role AS ENUM ('admin', 'investor', 'borrower', 'field_validator','field_officer');

CREATE TABLE IF NOT EXISTS "user".users (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email varchar unique,
  role "user".user_role not null,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS file.files (
  file_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  label varchar,
  location varchar NOT NULL,
  location_type varchar NOT NULL,
  file_type varchar,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp 
);

CREATE TABLE IF NOT EXISTS loan.loans (
  loan_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "description" varchar,
  proposed_by UUID not null REFERENCES "user".users (user_id),
  amount integer,
  duration_month integer,
  rate float,
  "state" loan.loan_state not null default 'proposed',
  approval_date timestamp,
  aggrement_file UUID REFERENCES file.files (file_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS loan.approvals (
  approval_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  loan_id UUID not null,
  visited_file UUID REFERENCES file.files (file_id),
  approved_by UUID not null REFERENCES "user".users (user_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS loan.investments (
  investment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  loan_id UUID not null,
  amount integer,
  invested_by UUID not null REFERENCES "user".users (user_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS loan.disbursments (
  disbursment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  loan_id UUID not null,
  disbursment_date timestamp,
  disbursment_by UUID not null REFERENCES "user".users (user_id),
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

ALTER TABLE loan.investments ADD CONSTRAINT check_total_investment
CHECK 
    ((SELECT SUM(amount) FROM loan.investment WHERE loan_id = loan.investment.loan_id) 
    <= (SELECT amount FROM loan.loans WHERE loan_id = loan.investment.loan_id));
