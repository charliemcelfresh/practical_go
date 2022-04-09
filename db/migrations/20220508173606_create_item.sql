-- migrate:up

CREATE TABLE item (
    item_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at timestamptz DEFAULT NOW()
);

-- migrate:down

DROP TABLE item;