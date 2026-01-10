package api

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type updateOrgRequest struct {
	Name string `json:"name"`
}

func handleGetOrg(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		var name string
		if err := db.QueryRow(`SELECT name FROM organizations WHERE id = ?`, orgID).Scan(&name); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "org not found")
		}
		return c.JSON(fiber.Map{"id": orgID, "name": name})
	}
}

func handleUpdateOrg(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := orgIDFrom(c)
		var req updateOrgRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		name := strings.TrimSpace(req.Name)
		if name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name required")
		}
		if _, err := db.Exec(`UPDATE organizations SET name = ? WHERE id = ?`, name, orgID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		return c.JSON(fiber.Map{"id": orgID, "name": name})
	}
}
