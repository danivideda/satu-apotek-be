ALTER TABLE IF EXISTS owners
DROP CONSTRAINT IF EXISTS owners_email_check_not_empty,
DROP CONSTRAINT IF EXISTS owners_username_check_not_empty;