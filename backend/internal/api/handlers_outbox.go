package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func handleListOutbox(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		rows, err := db.Query(`SELECT id, reminder_id, to_email, subject, body, created_at FROM outbox WHERE org_id = ? ORDER BY created_at DESC`, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer rows.Close()

		items := make([]fiber.Map, 0)
		for rows.Next() {
			var id, reminderID, toEmail, subject, body, createdAt string
			if err := rows.Scan(&id, &reminderID, &toEmail, &subject, &body, &createdAt); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			items = append(items, fiber.Map{
				"id": id,
				"reminder_id": reminderID,
				"to_email": toEmail,
				"subject": subject,
				"body": body,
				"created_at": createdAt,
			})
		}
		return c.JSON(fiber.Map{"outbox": items})
	}
}
