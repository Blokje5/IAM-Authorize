CREATE TABLE IF NOT EXISTS policies(
    id SERIAL PRIMARY KEY,
    version VARCHAR (3) NOT NULL,
    created_by VARCHAR (300) NOT NULL,
    last_modified_by VARCHAR (300) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_modified_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TYPE effect AS ENUM ('allow', 'deny');

CREATE TABLE IF NOT EXISTS statements(
    id SERIAL PRIMARY KEY,
    policy_id SERIAL REFERENCES policy(id) ON DELETE CASCADE
    effect effect NOT NULL,
    actions TEXT[],
    resources TEXT[],
);

CREATE OR REPLACE FUNCTION insert_policy(_version VARCHAR (3), _created_by VARCHAR (300), _ast_modified_by VARCHAR (300), _created_at TIMESTAMP, _last_modified_at TIMESTAMP)
RETURNS void
AS $$
    INSERT INTO policies (version, created_by, last_modified_by, created_at, last_modified_at)
    VALUES (_version, _created_by, _last_modified_by, _created_at, _last_modified_at)
$$
LANGUAGE sql;

CREATE OR REPLACE FUNCTION insert_statement(_policy_id UUID, _effect effect, _actions TEXT[], _resources TEXT[])
RETURNS void
AS $$
    INSERT INTO statements (policy_id, effect, actions, resources)
    VALUES (_policy_id, _effect, _actions, _resources)
$$
LANGUAGE sql;

CREATE OR REPLACE FUNCTION query_policy(_id SERIAL)
RETURNS json
AS $$
    WITH data AS (
        SELECT
            p.id,
            p.version,
            p.created_by,
            p.last_modified_by,
            p.created_at,
            p.last_modified_at,
            COALESCE(json_agg(s) FILTER (WHERE s IS NOT NULL), '[]') as statements
        FROM policies p
        LEFT JOIN statements s ON s.policy_id = p.id
        WHERE p.id = _id
        GROUP BY p.id
    )
    SELECT row_to_json(data) AS policy FROM data
$$
LANGUAGE sql;

CREATE OR REPLACE FUNCTION query_policies()
RETURNS setof json
AS $$
    WITH data AS (
        SELECT
            p.id,
            p.version,
            p.created_by,
            p.last_modified_by,
            p.created_at,
            p.last_modified_at,
            COALESCE(json_agg(s) FILTER (WHERE s IS NOT NULL), '[]') as statements
        FROM policies p
        LEFT JOIN statements s ON s.policy_id = p.id
        GROUP BY p.id
    )
    SELECT row_to_json(data) AS policy FROM data
$$
LANGUAGE sql;