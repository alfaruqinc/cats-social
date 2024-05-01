BEGIN;

CREATE TABLE IF NOT EXISTS cat_matches (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULl,
    issued_by UUID NOT NULL,
    match_cat_id UUID NOT NULL,
    user_cat_id UUID NOT NULL,
    message VARCHAR(120) NOT NULL
);

ALTER TABLE cat_matches ADD CONSTRAINT fk_issued_by_users FOREIGN KEY (issued_by) REFERENCES users (id);

ALTER TABLE cat_matches ADD CONSTRAINT fk_match_cat_id_cats FOREIGN KEY (match_cat_id) REFERENCES cats (id);

ALTER TABLE cat_matches ADD CONSTRAINT fk_user_cat_id_cats FOREIGN KEY (user_cat_id) REFERENCES cats (id);

COMMIT;
