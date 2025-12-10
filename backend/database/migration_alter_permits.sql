-- Migration ALTER for Permits table
-- This migrates existing permits table from old structure to new structure
-- Run this AFTER the initial migration if you already have data

-- Step 1: Add new columns
ALTER TABLE permits ADD COLUMN IF NOT EXISTS name VARCHAR(255);
ALTER TABLE permits ADD COLUMN IF NOT EXISTS responsible_person_id BIGINT;
ALTER TABLE permits ADD COLUMN IF NOT EXISTS responsible_doc_person_id BIGINT;

-- Step 2: Migrate data from old columns to new columns
UPDATE permits SET name = equipment_name WHERE name IS NULL;
-- Note: responsible_person (string) cannot be automatically migrated to responsible_person_id (foreign key)
-- You need to manually map the names to user IDs or set them to NULL

-- Step 3: Add foreign key constraints
ALTER TABLE permits ADD CONSTRAINT fk_permits_responsible_person 
    FOREIGN KEY (responsible_person_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE permits ADD CONSTRAINT fk_permits_responsible_doc_person 
    FOREIGN KEY (responsible_doc_person_id) REFERENCES users(id) ON DELETE SET NULL;

-- Step 4: Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_permits_name ON permits(name);
CREATE INDEX IF NOT EXISTS idx_permits_responsible_person_id ON permits(responsible_person_id);
CREATE INDEX IF NOT EXISTS idx_permits_responsible_doc_person_id ON permits(responsible_doc_person_id);

-- Step 5: Drop old index
DROP INDEX IF EXISTS idx_permits_equipment_name;

-- Step 6: Make name NOT NULL after data migration
ALTER TABLE permits ALTER COLUMN name SET NOT NULL;

-- Step 7: Drop old columns (CAREFUL: This will delete data!)
-- Uncomment these lines only after you've confirmed the migration is successful
-- ALTER TABLE permits DROP COLUMN IF EXISTS equipment_name;
-- ALTER TABLE permits DROP COLUMN IF EXISTS responsible_person;

-- Step 8: Update comments
COMMENT ON COLUMN permits.name IS 'Name of permit (equipment/service/competency/operational name)';
COMMENT ON COLUMN permits.responsible_person_id IS 'Reference to user responsible for this permit';
COMMENT ON COLUMN permits.responsible_doc_person_id IS 'Reference to user responsible for permit documentation';
