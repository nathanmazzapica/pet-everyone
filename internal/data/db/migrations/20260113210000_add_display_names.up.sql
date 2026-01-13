-- Add global display names and backfill existing users/guests

CREATE TYPE displayname_owner AS ENUM ('user', 'guest');

CREATE TABLE DisplayName (
    name VARCHAR(20) PRIMARY KEY,
    owner_type displayname_owner NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT display_name_owner_unique UNIQUE (owner_type, owner_id)
);

CREATE UNIQUE INDEX displayname_name_lower_idx ON DisplayName (lower(name));

ALTER TABLE RegisteredUser ADD COLUMN display_name VARCHAR(20);
ALTER TABLE Visitor ADD COLUMN display_name VARCHAR(20);

CREATE OR REPLACE FUNCTION generate_display_name()
RETURNS text AS $$
DECLARE
    adjectives text[] := ARRAY['brave','calm','swift','bright','happy','kind','quick','sharp','sunny','lucky'];
    nouns text[] := ARRAY['Fox','Cat','Dog','Hawk','Panda','Goat','Wolf','Lynx','Seal','Crab'];
    adjective text;
    noun text;
    num text;
    candidate text;
BEGIN
    LOOP
        adjective := adjectives[1 + floor(random() * array_length(adjectives, 1))::int];
        noun := nouns[1 + floor(random() * array_length(nouns, 1))::int];
        num := lpad((floor(random() * 10000))::text, 4, '0');
        candidate := adjective || noun || num;
        IF length(candidate) <= 20 THEN
            RETURN candidate;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION assign_display_name(owner_type displayname_owner, owner uuid)
RETURNS text AS $$
DECLARE
    candidate text;
BEGIN
    LOOP
        candidate := generate_display_name();
        BEGIN
            INSERT INTO DisplayName(name, owner_type, owner_id) VALUES (candidate, owner_type, owner);
            RETURN candidate;
        EXCEPTION WHEN unique_violation THEN
            -- retry on collision
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

UPDATE RegisteredUser
SET display_name = assign_display_name('user', user_id)
WHERE display_name IS NULL;

UPDATE Visitor
SET display_name = assign_display_name('guest', guest_id)
WHERE display_name IS NULL;

ALTER TABLE RegisteredUser ALTER COLUMN display_name SET NOT NULL;
ALTER TABLE Visitor ALTER COLUMN display_name SET NOT NULL;

ALTER TABLE RegisteredUser
    ADD CONSTRAINT registereduser_display_name_fk FOREIGN KEY (display_name) REFERENCES DisplayName(name) ON UPDATE CASCADE;
ALTER TABLE Visitor
    ADD CONSTRAINT visitor_display_name_fk FOREIGN KEY (display_name) REFERENCES DisplayName(name) ON UPDATE CASCADE;

DROP FUNCTION assign_display_name(displayname_owner, uuid);
DROP FUNCTION generate_display_name();
