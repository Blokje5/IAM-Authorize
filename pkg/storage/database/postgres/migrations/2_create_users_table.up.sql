CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR (300) UNIQUE NOT NULL,
    created_by VARCHAR (300) NOT NULL,
    last_modified_by VARCHAR (300) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_modified_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users_policies (
    user_id INTEGER REFERENCES users(id) ON DELETE NO ACTION,
    policy_id INTEGER REFERENCES policies(id) ON DELETE NO ACTION,
    PRIMARY KEY (user_id, policy_id)
);

CREATE OR REPLACE FUNCTION insert_user(_name VARCHAR (300), _created_by VARCHAR (300), _last_modified_by VARCHAR (300), _created_at TIMESTAMP, _last_modified_at TIMESTAMP)
RETURNS INTEGER
AS $$
    INSERT INTO users (name, created_by, last_modified_by, created_at, last_modified_at)
    VALUES (_name, _created_by, _last_modified_by, _created_at, _last_modified_at)
    RETURNING ID
$$
LANGUAGE sql;

CREATE OR REPLACE FUNCTION query_policies_for_user(_name VARCHAR (300))
RETURNS setof json
AS $$
    WITH user AS (
        SELECT id
        FROM users
        WHERE name = _name
    ), user_policies AS (
        SELECT policy_id AS id
        FROM users_policies
        INNER JOIN user ON user.id = users_policies.user_id
    ), data AS (
        SELECT
            p.id,
            p.version,
            p.created_by,
            p.last_modified_by,
            p.created_at,
            p.last_modified_at,
            COALESCE(json_agg(s) FILTER (WHERE s IS NOT NULL), '[]') as statements
        FROM policies p
        INNER JOIN user_policies up ON p.id = up.id
        LEFT JOIN statements s ON s.policy_id = p.id
        GROUP BY p.id
    )
    SELECT row_to_json(data) AS policy FROM data
$$
LANGUAGE sql;