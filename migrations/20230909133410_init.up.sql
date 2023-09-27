BEGIN;

CREATE TABLE white_list (
    id SERIAL PRIMARY KEY,
    ip_address CIDR NOT NULL
);
CREATE INDEX idx_whitelist_ip_address ON white_list USING GIST (ip_address inet_ops);

CREATE TABLE black_list (
    id SERIAL PRIMARY KEY,
    ip_address CIDR NOT NULL
);
CREATE INDEX idx_blacklist_ip_address ON black_list USING GIST (ip_address inet_ops);

COMMIT;
