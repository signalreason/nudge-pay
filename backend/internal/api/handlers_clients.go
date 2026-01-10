package api

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type clientPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
	Notes   string `json:"notes"`
}

func handleListClients(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		rows, err := db.Query(`SELECT id, name, email, company, phone, notes, created_at FROM clients WHERE org_id = ? ORDER BY created_at DESC`, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer rows.Close()

		clients := make([]fiber.Map, 0)
		for rows.Next() {
			var id, name, email, company, phone, notes, createdAt string
			if err := rows.Scan(&id, &name, &email, &company, &phone, &notes, &createdAt); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "db error")
			}
			clients = append(clients, fiber.Map{
				"id": id, "name": name, "email": email, "company": company, "phone": phone, "notes": notes, "created_at": createdAt,
			})
		}
		return c.JSON(fiber.Map{"clients": clients})
	}
}

func handleCreateClient(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		var req clientPayload
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Name == "" || req.Email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name and email required")
		}
		if req.Company == "" {
			req.Company = "-"
		}
		id := uuid.NewString()
		now := time.Now().UTC().Format(time.RFC3339)
		if _, err := db.Exec(`INSERT INTO clients (id, org_id, name, email, company, phone, notes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			id, orgID, req.Name, req.Email, req.Company, req.Phone, req.Notes, now); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
	}
}

func handleGetClient(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		var name, email, company, phone, notes, createdAt string
		if err := db.QueryRow(`SELECT name, email, company, phone, notes, created_at FROM clients WHERE id = ? AND org_id = ?`, id, orgID).
			Scan(&name, &email, &company, &phone, &notes, &createdAt); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "client not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		return c.JSON(fiber.Map{
			"id": id, "name": name, "email": email, "company": company, "phone": phone, "notes": notes, "created_at": createdAt,
		})
	}
}

func handleUpdateClient(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		var req clientPayload
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		name := strings.TrimSpace(req.Name)
		email := strings.TrimSpace(strings.ToLower(req.Email))
		if name == "" || email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name and email required")
		}
		if req.Company == "" {
			req.Company = "-"
		}
		res, err := db.Exec(`UPDATE clients SET name = ?, email = ?, company = ?, phone = ?, notes = ? WHERE id = ? AND org_id = ?`,
			name, email, req.Company, req.Phone, req.Notes, id, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "client not found")
		}
		return c.JSON(fiber.Map{"id": id})
	}
}

func handleDeleteClient(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		id := c.Params("id")
		res, err := db.Exec(`DELETE FROM clients WHERE id = ? AND org_id = ?`, id, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "client not found")
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
