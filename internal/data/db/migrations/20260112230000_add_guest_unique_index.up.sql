CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS uniq_guest_pet
    ON UserPetsClickCount (pet_id, guest_id)
    WHERE guest_id IS NOT NULL;
