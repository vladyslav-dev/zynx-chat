-- Step 1: Alter the email column to be NOT NULL
ALTER TABLE users
ALTER COLUMN email SET NOT NULL;

-- Step 2: Remove the phone column
ALTER TABLE users
DROP COLUMN phone;