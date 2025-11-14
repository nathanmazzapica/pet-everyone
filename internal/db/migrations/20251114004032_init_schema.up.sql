-- Schema UP
CREATE TABLE Visitor (
    guest_id UUID PRIMARY KEY NOT NULL,
    last_seen TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE RegisteredUser (
    user_id UUID PRIMARY KEY NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE PetImage (
    image_id SERIAL PRIMARY KEY NOT NULL,
    image_url TEXT NOT NULL DEFAULT '/assets/placeholder.jpg',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE Pet (
    pet_id UUID PRIMARY KEY NOT NULL,
    pet_name VARCHAR(128) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    visibility BOOLEAN NOT NULL DEFAULT TRUE,
    user_id UUID,
    active_image SERIAL,

    FOREIGN KEY (user_id) REFERENCES RegisteredUser(user_id),
    FOREIGN KEY (active_image) REFERENCES PetImage(image_id)
);

CREATE TABLE UserPetsClickCount
(
    click_count INT default 0,
    pet_id      UUID NOT NULL,
    user_id     UUID,
    guest_id    UUID,

    FOREIGN KEY (pet_id) REFERENCES Pet (pet_id),
    FOREIGN KEY (user_id) REFERENCES RegisteredUser (user_id),
    FOREIGN KEY (guest_id) REFERENCES Visitor (guest_id),

    CHECK (
        (user_id IS NOT NULL AND guest_id IS NULL) OR
        (user_id IS NULL AND guest_id IS NOT NULL)
    ),

    PRIMARY KEY (pet_id, user_id, guest_id)
);