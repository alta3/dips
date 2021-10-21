PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE host (
    mac INTEGER NOT NULL,
    ip INTEGER NOT NULL,
    hostname TEXT NOT NULL,
    domain TEXT NOT NULL,
    gateway INTEGER NOT NULL,
    network TEXT NOT NULL,
    requestor TEXT NOT NULL,

    UNIQUE(mac),
    UNIQUE(ip),
    UNIQUE(hostname, domain)
    --FOREIGN KEY (requestor) REFERENCES requestor(name)
);
CREATE TABLE requestor (
    name TEXT PRIMARY KEY,
    api_key TEXT NOT NULL,
    UNIQUE(api_key)
);
COMMIT;
