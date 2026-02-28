package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type ProdOrg struct {
	ID        int64
	Alias     string
	IsActive  bool
	CreatedAt time.Time
}

type Sandbox struct {
	ID        int64
	Name      string
	ProdOrgID int64
	CreatedAt time.Time
}

type Store interface {
	AddProdOrg(alias string) error
	RemoveProdOrg(id int64) error
	ListProdOrgs() ([]ProdOrg, error)
	SetActiveProdOrg(id int64) error
	GetActiveProdOrg() (*ProdOrg, error)

	AddSandbox(name string, prodOrgID int64) error
	EnsureSandbox(name string, prodOrgID int64) error
	RemoveSandbox(id int64) error
	RemoveSandboxByName(name string, prodOrgID int64) error
	ListSandboxes(prodOrgID int64) ([]Sandbox, error)

	Close() error
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) AddProdOrg(alias string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// If this is the first org, make it active
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM prod_orgs").Scan(&count); err != nil {
		return fmt.Errorf("count prod orgs: %w", err)
	}

	isActive := 0
	if count == 0 {
		isActive = 1
	}

	_, err = tx.Exec("INSERT INTO prod_orgs (alias, is_active) VALUES (?, ?)", alias, isActive)
	if err != nil {
		return fmt.Errorf("insert prod org: %w", err)
	}

	return tx.Commit()
}

func (s *SQLiteStore) RemoveProdOrg(id int64) error {
	_, err := s.db.Exec("DELETE FROM prod_orgs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete prod org: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListProdOrgs() ([]ProdOrg, error) {
	rows, err := s.db.Query("SELECT id, alias, is_active, created_at FROM prod_orgs ORDER BY alias")
	if err != nil {
		return nil, fmt.Errorf("query prod orgs: %w", err)
	}
	defer rows.Close()

	var orgs []ProdOrg
	for rows.Next() {
		var o ProdOrg
		var createdAt string
		if err := rows.Scan(&o.ID, &o.Alias, &o.IsActive, &createdAt); err != nil {
			return nil, fmt.Errorf("scan prod org: %w", err)
		}
		o.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		orgs = append(orgs, o)
	}
	return orgs, rows.Err()
}

func (s *SQLiteStore) SetActiveProdOrg(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE prod_orgs SET is_active = 0"); err != nil {
		return fmt.Errorf("deactivate all: %w", err)
	}
	if _, err := tx.Exec("UPDATE prod_orgs SET is_active = 1 WHERE id = ?", id); err != nil {
		return fmt.Errorf("activate prod org: %w", err)
	}

	return tx.Commit()
}

func (s *SQLiteStore) GetActiveProdOrg() (*ProdOrg, error) {
	var o ProdOrg
	var createdAt string
	err := s.db.QueryRow(
		"SELECT id, alias, is_active, created_at FROM prod_orgs WHERE is_active = 1",
	).Scan(&o.ID, &o.Alias, &o.IsActive, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query active prod org: %w", err)
	}
	o.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &o, nil
}

func (s *SQLiteStore) AddSandbox(name string, prodOrgID int64) error {
	_, err := s.db.Exec(
		"INSERT INTO sandboxes (name, prod_org_id) VALUES (?, ?)",
		name, prodOrgID,
	)
	if err != nil {
		return fmt.Errorf("insert sandbox: %w", err)
	}
	return nil
}

func (s *SQLiteStore) EnsureSandbox(name string, prodOrgID int64) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO sandboxes (name, prod_org_id) VALUES (?, ?)",
		name, prodOrgID,
	)
	if err != nil {
		return fmt.Errorf("ensure sandbox: %w", err)
	}
	return nil
}

func (s *SQLiteStore) RemoveSandbox(id int64) error {
	_, err := s.db.Exec("DELETE FROM sandboxes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete sandbox: %w", err)
	}
	return nil
}

func (s *SQLiteStore) RemoveSandboxByName(name string, prodOrgID int64) error {
	_, err := s.db.Exec("DELETE FROM sandboxes WHERE name = ? AND prod_org_id = ?", name, prodOrgID)
	if err != nil {
		return fmt.Errorf("delete sandbox by name: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListSandboxes(prodOrgID int64) ([]Sandbox, error) {
	rows, err := s.db.Query(
		"SELECT id, name, prod_org_id, created_at FROM sandboxes WHERE prod_org_id = ? ORDER BY name",
		prodOrgID,
	)
	if err != nil {
		return nil, fmt.Errorf("query sandboxes: %w", err)
	}
	defer rows.Close()

	var sandboxes []Sandbox
	for rows.Next() {
		var sb Sandbox
		var createdAt string
		if err := rows.Scan(&sb.ID, &sb.Name, &sb.ProdOrgID, &createdAt); err != nil {
			return nil, fmt.Errorf("scan sandbox: %w", err)
		}
		sb.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		sandboxes = append(sandboxes, sb)
	}
	return sandboxes, rows.Err()
}
