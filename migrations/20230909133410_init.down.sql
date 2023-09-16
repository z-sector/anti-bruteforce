BEGIN;

DROP INDEX IF EXISTS idx_whitelist_ip_address;
DROP INDEX IF EXISTS idx_blacklist_ip_address;

DROP TABLE IF EXISTS whitelist;
DROP TABLE IF EXISTS blacklist;

COMMIT;