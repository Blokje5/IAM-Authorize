CREATE TABLE IF NOT EXISTS namespaces(
    id SERIAL PRIMARY KEY,
    name VARCHAR (50) UNIQUE NOT NULL,
    created_by VARCHAR (300) NOT NULL,
    last_modified_by VARCHAR (300) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_modified_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);