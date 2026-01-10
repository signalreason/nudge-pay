package api

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nudgepay/internal/services"
)

type invoicePayload struct {
	ClientID        string `json:"client_id"`
	TemplateID      string `json:"template_id"`
	Number          string `json:"number"`
	AmountCents     int64  `json:"amount_cents"`
	Currency        string `json:"currency"`
	DueDate         string `json:"due_date"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
	ReminderOffsets []int  `json:"reminder_offsets"`
}

func handleListInvoices(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		status := strings.TrimSpace(c.Query("status"))
		query := `SELECT id, client_id, template_id, number, amount_cents, currency, due_date, status, notes, created_at, updated_at
			FROM invoices WHERE org_id = ?`
		args := []interface{}{orgID}
		if status != "" {
			query += " AND status = ?"
			args = append(args, status)
		}
		query += " ORDER BY created_at DESC"
		rows, err := db.Query(query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer rows.Close()

		invoices := make([]fiber.Map, 0)
		for rows.Next() {
			var id, clientID, number, currency, dueDate, status, notes, createdAt, updatedAt string
			var templateID sql.NullString
			var amountCents int64
			if err := rows.Scan(&id, &clientID, &templateID, &number, &amountCents, &currency, &dueDate, &status, &notes, &createdAt, &updatedAt); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			invoices = append(invoices, fiber.Map{
				"id": id,
				"client_id": clientID,
				"template_id": nullIfEmpty(templateID.String),
				"number": number,
				"amount_cents": amountCents,
				"currency": currency,
				"due_date": dueDate,
				"status": status,
				"notes": notes,
				"created_at": createdAt,
				"updated_at": updatedAt,
			})
		}
		return c.JSON(fiber.Map{"invoices": invoices})
	}
}

func handleCreateInvoice(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		var req invoicePayload
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		req.ClientID = strings.TrimSpace(req.ClientID)
		req.Number = strings.TrimSpace(req.Number)
		req.Currency = strings.ToUpper(strings.TrimSpace(req.Currency))
		if req.ClientID == "" || req.Number == "" || req.AmountCents <= 0 || req.Currency == "" || req.DueDate == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing required fields")
		}

		if req.Status == "" {
			req.Status = "sent"
		}
		if req.Notes == "" {
			req.Notes = ""
		}
		if len(req.ReminderOffsets) == 0 {
			req.ReminderOffsets = []int{-3, 0, 7}
		}

		var clientExists string
		if err := db.QueryRow(`SELECT id FROM clients WHERE id = ? AND org_id = ?`, req.ClientID, orgID).Scan(&clientExists); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "client not found")
		}

		if req.TemplateID != "" {
			var tmp string
			if err := db.QueryRow(`SELECT id FROM templates WHERE id = ? AND org_id = ?`, req.TemplateID, orgID).Scan(&tmp); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "template not found")
			}
		} else {
			tid, err := services.EnsureDefaultTemplate(db, orgID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			req.TemplateID = tid
		}

		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid due_date")
		}
		dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, time.UTC)

		invoiceID := uuid.NewString()
		now := time.Now().UTC().Format(time.RFC3339)
		if _, err := db.Exec(`INSERT INTO invoices (id, org_id, client_id, template_id, number, amount_cents, currency, due_date, status, notes, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			invoiceID, orgID, req.ClientID, req.TemplateID, req.Number, req.AmountCents, req.Currency, dueDate.Format(time.RFC3339), req.Status, req.Notes, now, now); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		for _, offset := range req.ReminderOffsets {
			scheduled := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 9, 0, 0, 0, time.UTC).AddDate(0, 0, offset)
			reminderID := uuid.NewString()
			if _, err := db.Exec(`INSERT INTO reminders (id, org_id, invoice_id, template_id, scheduled_for, sent_at, status, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				reminderID, orgID, invoiceID, req.TemplateID, scheduled.Format(time.RFC3339), nil, "scheduled", now); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": invoiceID})
	}
}

func handleGetInvoice(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		var clientID, number, currency, dueDate, status, notes, createdAt, updatedAt string
		var templateID sql.NullString
		var amountCents int64
		if err := db.QueryRow(`SELECT client_id, template_id, number, amount_cents, currency, due_date, status, notes, created_at, updated_at
			FROM invoices WHERE id = ? AND org_id = ?`, id, orgID).
			Scan(&clientID, &templateID, &number, &amountCents, &currency, &dueDate, &status, &notes, &createdAt, &updatedAt); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "invoice not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		rows, err := db.Query(`SELECT id, scheduled_for, sent_at, status FROM reminders WHERE invoice_id = ? ORDER BY scheduled_for ASC`, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer rows.Close()
		reminders := make([]fiber.Map, 0)
		for rows.Next() {
			var reminderID, scheduledFor, rStatus string
			var sentAt sql.NullString
			if err := rows.Scan(&reminderID, &scheduledFor, &sentAt, &rStatus); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			reminders = append(reminders, fiber.Map{
				"id": reminderID, "scheduled_for": scheduledFor, "sent_at": nullIfEmpty(sentAt.String), "status": rStatus,
			})
		}

		return c.JSON(fiber.Map{
			"id": id,
			"client_id": clientID,
			"template_id": nullIfEmpty(templateID.String),
			"number": number,
			"amount_cents": amountCents,
			"currency": currency,
			"due_date": dueDate,
			"status": status,
			"notes": notes,
			"created_at": createdAt,
			"updated_at": updatedAt,
			"reminders": reminders,
		})
	}
}

func handleUpdateInvoice(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		var req invoicePayload
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		if req.AmountCents < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "amount_cents invalid")
		}
		if req.Currency != "" {
			req.Currency = strings.ToUpper(strings.TrimSpace(req.Currency))
		}
		if req.DueDate != "" {
			if _, err := time.Parse("2006-01-02", req.DueDate); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid due_date")
			}
		}
		if req.TemplateID != "" {
			var tmp string
			if err := db.QueryRow(`SELECT id FROM templates WHERE id = ? AND org_id = ?`, req.TemplateID, orgID).Scan(&tmp); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "template not found")
			}
		}

		fields := []string{}
		args := []interface{}{}
		if req.Number != "" {
			fields = append(fields, "number = ?")
			args = append(args, strings.TrimSpace(req.Number))
		}
		if req.AmountCents > 0 {
			fields = append(fields, "amount_cents = ?")
			args = append(args, req.AmountCents)
		}
		if req.Currency != "" {
			fields = append(fields, "currency = ?")
			args = append(args, req.Currency)
		}
		if req.DueDate != "" {
			parsed, _ := time.Parse("2006-01-02", req.DueDate)
			parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
			fields = append(fields, "due_date = ?")
			args = append(args, parsed.Format(time.RFC3339))
		}
		if req.Status != "" {
			fields = append(fields, "status = ?")
			args = append(args, strings.TrimSpace(req.Status))
		}
		if req.Notes != "" {
			fields = append(fields, "notes = ?")
			args = append(args, req.Notes)
		}
		if req.TemplateID != "" {
			fields = append(fields, "template_id = ?")
			args = append(args, req.TemplateID)
		}

		if len(fields) == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
		}
		fields = append(fields, "updated_at = ?")
		args = append(args, time.Now().UTC().Format(time.RFC3339))
		args = append(args, id, orgID)

		query := `UPDATE invoices SET ` + strings.Join(fields, ", ") + ` WHERE id = ? AND org_id = ?`
		res, err := db.Exec(query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "invoice not found")
		}
		return c.JSON(fiber.Map{"id": id})
	}
}

func handleDeleteInvoice(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		res, err := db.Exec(`DELETE FROM invoices WHERE id = ? AND org_id = ?`, id, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "invoice not found")
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func nullIfEmpty(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
