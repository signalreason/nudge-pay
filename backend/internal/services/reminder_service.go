package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReminderInfo struct {
	ID         string
	InvoiceID  string
	TemplateID string
}

func SendDueReminders(db *sql.DB, orgID string, now time.Time) (int, error) {
	rows, err := db.Query(`SELECT id, invoice_id, template_id FROM reminders
		WHERE org_id = ? AND status = 'scheduled' AND scheduled_for <= ?`, orgID, now.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	reminders := make([]ReminderInfo, 0)
	for rows.Next() {
		var info ReminderInfo
		var templateID sql.NullString
		if err := rows.Scan(&info.ID, &info.InvoiceID, &templateID); err != nil {
			return 0, err
		}
		info.TemplateID = templateID.String
		reminders = append(reminders, info)
	}

	sent := 0
	for _, reminder := range reminders {
		ok, err := sendReminder(db, orgID, reminder.ID, reminder.InvoiceID, reminder.TemplateID, now)
		if err != nil {
			return sent, err
		}
		if ok {
			sent++
		}
	}
	return sent, nil
}

func SendReminderByID(db *sql.DB, orgID, reminderID string, now time.Time) (bool, error) {
	var invoiceID string
	var templateID sql.NullString
	if err := db.QueryRow(`SELECT invoice_id, template_id FROM reminders WHERE id = ? AND org_id = ?`, reminderID, orgID).
		Scan(&invoiceID, &templateID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return sendReminder(db, orgID, reminderID, invoiceID, templateID.String, now)
}

func sendReminder(db *sql.DB, orgID, reminderID, invoiceID, templateID string, now time.Time) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`UPDATE reminders SET status = 'sent', sent_at = ? WHERE id = ? AND org_id = ? AND status = 'scheduled'`,
		now.Format(time.RFC3339), reminderID, orgID)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return false, nil
	}

	var clientName, clientEmail, clientCompany string
	var invoiceNumber, currency, dueDate string
	var amountCents int64
	if err := tx.QueryRow(`SELECT c.name, c.email, c.company, i.number, i.amount_cents, i.currency, i.due_date
		FROM invoices i JOIN clients c ON i.client_id = c.id
		WHERE i.id = ? AND i.org_id = ?`, invoiceID, orgID).
		Scan(&clientName, &clientEmail, &clientCompany, &invoiceNumber, &amountCents, &currency, &dueDate); err != nil {
		return false, err
	}

	if strings.TrimSpace(templateID) == "" {
		defaultID, err := EnsureDefaultTemplate(db, orgID)
		if err != nil {
			return false, err
		}
		templateID = defaultID
	}

	var subject, body string
	if err := tx.QueryRow(`SELECT subject, body FROM templates WHERE id = ? AND org_id = ?`, templateID, orgID).
		Scan(&subject, &body); err != nil {
		if err == sql.ErrNoRows {
			fallbackID, err := EnsureDefaultTemplate(db, orgID)
			if err != nil {
				return false, err
			}
			if err := tx.QueryRow(`SELECT subject, body FROM templates WHERE id = ? AND org_id = ?`, fallbackID, orgID).
				Scan(&subject, &body); err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	var orgName string
	if err := tx.QueryRow(`SELECT name FROM organizations WHERE id = ?`, orgID).Scan(&orgName); err != nil {
		return false, err
	}

	amount := formatAmount(amountCents, currency)
	values := map[string]string{
		"client_name":   clientName,
		"client_company": clientCompany,
		"invoice_number": invoiceNumber,
		"amount":         amount,
		"due_date":       dueDate,
		"org_name":       orgName,
	}

	finalSubject := applyTemplate(subject, values)
	finalBody := applyTemplate(body, values)

	outboxID := uuid.NewString()
	if _, err := tx.Exec(`INSERT INTO outbox (id, org_id, reminder_id, to_email, subject, body, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		outboxID, orgID, reminderID, clientEmail, finalSubject, finalBody, now.Format(time.RFC3339)); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func applyTemplate(input string, values map[string]string) string {
	out := input
	for key, value := range values {
		out = strings.ReplaceAll(out, fmt.Sprintf("{{%s}}", key), value)
	}
	return out
}

func formatAmount(amountCents int64, currency string) string {
	amount := float64(amountCents) / 100.0
	return fmt.Sprintf("%s %.2f", strings.ToUpper(currency), amount)
}
