--Initial Pets Up Migration
INSERT INTO Pet (pet_id, pet_name, visibility)
VALUES
(gen_random_uuid(), 'Henry', TRUE),
(gen_random_uuid(), 'Daisy', TRUE);