ALTER TABLE IF EXISTS owners
ADD CONSTRAINT owners_email_check_not_empty CHECK (email != ''),
ADD CONSTRAINT owners_username_check_not_empty CHECK (username != '');