package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func New(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			owner_user_id TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			org_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS clients (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			company TEXT NOT NULL,
			phone TEXT NOT NULL,
			notes TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS templates (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL,
			name TEXT NOT NULL,
			subject TEXT NOT NULL,
			body TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS invoices (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL,
			client_id TEXT NOT NULL,
			template_id TEXT,
			number TEXT NOT NULL,
			amount_cents INTEGER NOT NULL,
			currency TEXT NOT NULL,
			due_date TEXT NOT NULL,
			status TEXT NOT NULL,
			notes TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
			FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS reminders (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL,
			invoice_id TEXT NOT NULL,
			template_id TEXT,
			scheduled_for TEXT NOT NULL,
			sent_at TEXT,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
			FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE SET NULL
		);`,
		`CREATE TABLE IF NOT EXISTS outbox (
			id TEXT PRIMARY KEY,
			org_id TEXT NOT NULL,
			reminder_id TEXT NOT NULL,
			to_email TEXT NOT NULL,
			subject TEXT NOT NULL,
			body TEXT NOT NULL,
			created_at TEXT NOT NULL,
			FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE,
			FOREIGN KEY (reminder_id) REFERENCES reminders(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_clients_org ON clients(org_id);`,
		`CREATE INDEX IF NOT EXISTS idx_invoices_org ON invoices(org_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_due ON reminders(status, scheduled_for);`,
		`CREATE INDEX IF NOT EXISTS idx_outbox_org ON outbox(org_id);`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
