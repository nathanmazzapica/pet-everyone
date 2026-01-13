-- Remove global display names

ALTER TABLE Visitor DROP CONSTRAINT visitor_display_name_fk;
ALTER TABLE RegisteredUser DROP CONSTRAINT registereduser_display_name_fk;

ALTER TABLE Visitor DROP COLUMN display_name;
ALTER TABLE RegisteredUser DROP COLUMN display_name;

DROP INDEX IF EXISTS displayname_name_lower_idx;
DROP TABLE DisplayName;
DROP TYPE displayname_owner;
