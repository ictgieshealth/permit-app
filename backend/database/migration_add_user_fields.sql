-- Migration: Add phone_number and nip fields to users table
-- Created: 2025-12-09

-- Add phone_number column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20);

-- Add nip column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS nip VARCHAR(50);

-- Add comments to columns
COMMENT ON COLUMN users.phone_number IS 'User phone number';
COMMENT ON COLUMN users.nip IS 'Nomor Induk Pegawai (Employee ID Number)';
