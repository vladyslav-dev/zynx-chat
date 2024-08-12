-- Step 1: Add the phone column as nullable
ALTER TABLE users
ADD COLUMN phone VARCHAR(255);

-- Step 2: Update existing rows with a default value
UPDATE users
SET phone = '380000000000'; -- Replace 'default-value' with an appropriate value if needed

-- Step 3: Alter the column to enforce NOT NULL constraint
ALTER TABLE users
ALTER COLUMN phone SET NOT NULL;

-- Step 4: Alter the email column to be nullable
ALTER TABLE users
ALTER COLUMN email DROP NOT NULL;