--Initial Pets Up Migration
INSERT INTO PetImage (image_id, image_url)
VALUES
(1, '/assets/henry.png'),
(2, '/assets/DAISY.png');

INSERT INTO Pet (pet_id, pet_name, visibility)
VALUES
(gen_random_uuid(), 'Henry', TRUE),
(gen_random_uuid(), 'Daisy', TRUE);