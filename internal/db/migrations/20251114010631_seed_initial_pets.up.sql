--Initial Pets Up Migration
INSERT INTO PetImage (image_id, image_url)
VALUES
(1, '/assets/henry.png'),
(2, '/assets/DAISY.png');

INSERT INTO Pet (pet_id, pet_name, visibility)
VALUES
(gen_random_uuid(), 'Henry', TRUE),
(gen_random_uuid(), 'Daisy', TRUE);

INSERT INTO RegisteredUser (user_id, email, password_hash, created_at)
VALUES
(gen_random_uuid(), 'admin@pet.com', '$argon2id$v=19$m=65536,t=1,p=12$gb+w82wj5udmBIiq5x7QVw$BW9hN4LGq2OVY3Mizm5jrjo+E9Irhah3JihLAwj2wW0', now());