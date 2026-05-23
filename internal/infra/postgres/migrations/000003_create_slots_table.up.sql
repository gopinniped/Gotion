CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(30) NOT NULL
);

CREATE TABLE slots (
   id UUID PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   description TEXT,
   category_id UUID REFERENCES categories(id) ON DELETE RESTRICT,
   start_time TIMESTAMPTZ NOT NULL,
   end_time TIMESTAMPTZ NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);