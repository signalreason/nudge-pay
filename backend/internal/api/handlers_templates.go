package api

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type templatePayload struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func handleListTemplates(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		rows, err := db.Query(`SELECT id, name, subject, body, created_at, updated_at FROM templates WHERE org_id = ? ORDER BY created_at DESC`, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer rows.Close()

		templates := make([]fiber.Map, 0)
		for rows.Next() {
			var id, name, subject, body, createdAt, updatedAt string
			if err := rows.Scan(&id, &name, &subject, &body, &createdAt, &updatedAt); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			templates = append(templates, fiber.Map{
				"id": id, "name": name, "subject": subject, "body": body, "created_at": createdAt, "updated_at": updatedAt,
			})
		}
		return c.JSON(fiber.Map{"templates": templates})
	}
}

func handleCreateTemplate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		var req templatePayload
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		name := strings.TrimSpace(req.Name)
		subject := strings.TrimSpace(req.Subject)
		body := strings.TrimSpace(req.Body)
		if name == "" || subject == "" || body == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name, subject, and body required")
		}
		id := uuid.NewString()
		now := time.Now().UTC().Format(time.RFC3339)
		if _, err := db.Exec(`INSERT INTO templates (id, org_id, name, subject, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id, orgID, name, subject, body, now, now); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
	}
}

func handleUpdateTemplate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		var req templatePayload
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		name := strings.TrimSpace(req.Name)
		subject := strings.TrimSpace(req.Subject)
		body := strings.TrimSpace(req.Body)
		if name == "" || subject == "" || body == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name, subject, and body required")
		}
		now := time.Now().UTC().Format(time.RFC3339)
		res, err := db.Exec(`UPDATE templates SET name = ?, subject = ?, body = ?, updated_at = ? WHERE id = ? AND org_id = ?`,
			name, subject, body, now, id, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "template not found")
		}
		return c.JSON(fiber.Map{"id": id})
	}
}

func handleDeleteTemplate(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		res, err := db.Exec(`DELETE FROM templates WHERE id = ? AND org_id = ?`, id, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "template not found")
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
