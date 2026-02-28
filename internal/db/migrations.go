package db

import "database/sql"

const schema = `
CREATE TABLE IF NOT EXISTS prod_orgs (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    alias      TEXT    NOT NULL UNIQUE,
    is_active  INTEGER NOT NULL DEFAULT 0,
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS sandboxes (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL,
    prod_org_id INTEGER NOT NULL REFERENCES prod_orgs(id) ON DELETE CASCADE,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
    UNIQUE(name, prod_org_id)
);
`

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
