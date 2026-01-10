package services

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

const DefaultTemplateName = "Default Reminder"

func EnsureDefaultTemplate(db *sql.DB, orgID string) (string, error) {
	var id string
	err := db.QueryRow(`SELECT id FROM templates WHERE org_id = ? ORDER BY created_at ASC LIMIT 1`, orgID).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}

	id = uuid.NewString()
	now := time.Now().UTC().Format(time.RFC3339)
	subject := "Friendly reminder: invoice {{invoice_number}}"
	body := strings.Join([]string{
		"Hi {{client_name}},",
		"",
		"Just a quick reminder that invoice {{invoice_number}} for {{amount}} is due on {{due_date}}.",
		"If you've already sent payment, please disregard this note.",
		"",
		"Thanks,",
		"{{org_name}}",
	}, "\n")
	_, err = db.Exec(`INSERT INTO templates (id, org_id, name, subject, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, orgID, DefaultTemplateName, subject, body, now, now)
	if err != nil {
		return "", err
	}
	return id, nil
}
