package api

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"nudgepay/internal/services"
)

func handleListReminders(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		status := strings.TrimSpace(c.Query("status"))
		query := `SELECT r.id, r.invoice_id, r.scheduled_for, r.sent_at, r.status, i.number
			FROM reminders r JOIN invoices i ON r.invoice_id = i.id
			WHERE r.org_id = ?`
		args := []interface{}{orgID}
		if status != "" {
			query += " AND r.status = ?"
			args = append(args, status)
		}
		query += " ORDER BY r.scheduled_for ASC"

		rows, err := db.Query(query, args...)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer rows.Close()

		reminders := make([]fiber.Map, 0)
		for rows.Next() {
			var id, invoiceID, scheduledFor, status, number string
			var sentAt sql.NullString
			if err := rows.Scan(&id, &invoiceID, &scheduledFor, &sentAt, &status, &number); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			reminders = append(reminders, fiber.Map{
				"id": id,
				"invoice_id": invoiceID,
				"invoice_number": number,
				"scheduled_for": scheduledFor,
				"sent_at": nullIfEmpty(sentAt.String),
				"status": status,
			})
		}

		return c.JSON(fiber.Map{"reminders": reminders})
	}
}

func handleSendReminder(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		reminderID := c.Params("id")
		sent, err := services.SendReminderByID(db, orgID, reminderID, time.Now().UTC())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "send failed")
		}
		if !sent {
			return fiber.NewError(fiber.StatusNotFound, "reminder not found or already sent")
		}
		return c.JSON(fiber.Map{"id": reminderID, "status": "sent"})
	}
}

func handleSendDueReminders(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		sent, err := services.SendDueReminders(db, orgID, time.Now().UTC())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "send failed")
		}
		return c.JSON(fiber.Map{"sent": sent})
	}
}
