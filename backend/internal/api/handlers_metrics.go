package api

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

func handleMetrics(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)

		var clientCount int
		if err := db.QueryRow(`SELECT COUNT(*) FROM clients WHERE org_id = ?`, orgID).Scan(&clientCount); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		var invoiceCount int
		if err := db.QueryRow(`SELECT COUNT(*) FROM invoices WHERE org_id = ?`, orgID).Scan(&invoiceCount); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		now := time.Now().UTC()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
		var overdueCount int
		if err := db.QueryRow(`SELECT COUNT(*) FROM invoices WHERE org_id = ? AND status != 'paid' AND due_date < ?`, orgID, today).Scan(&overdueCount); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		windowEnd := now.Add(7 * 24 * time.Hour).Format(time.RFC3339)
		var upcomingReminders int
		if err := db.QueryRow(`SELECT COUNT(*) FROM reminders WHERE org_id = ? AND status = 'scheduled' AND scheduled_for BETWEEN ? AND ?`,
			orgID, now.Format(time.RFC3339), windowEnd).Scan(&upcomingReminders); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		var outstanding int64
		if err := db.QueryRow(`SELECT COALESCE(SUM(amount_cents), 0) FROM invoices WHERE org_id = ? AND status != 'paid'`, orgID).Scan(&outstanding); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		return c.JSON(fiber.Map{
			"clients": clientCount,
			"invoices": invoiceCount,
			"overdue": overdueCount,
			"upcoming_reminders": upcomingReminders,
			"outstanding_cents": outstanding,
		})
	}
}
