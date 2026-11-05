CREATE TABLE items (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(255) NOT NULL,
    description varchar(1024),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);
