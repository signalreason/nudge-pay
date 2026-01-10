package api

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"nudgepay/internal/auth"
	"nudgepay/internal/config"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OrgName  string `json:"org_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleRegister(db *sql.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req registerRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || req.Password == "" || strings.TrimSpace(req.OrgName) == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing required fields")
		}
		var exists string
		if err := db.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&exists); err == nil {
			return fiber.NewError(fiber.StatusConflict, "email already registered")
		} else if err != sql.ErrNoRows {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "password hashing failed")
		}

		orgID := uuid.NewString()
		userID := uuid.NewString()
		now := time.Now().UTC().Format(time.RFC3339)

		tx, err := db.Begin()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		defer tx.Rollback()

		if _, err := tx.Exec(`INSERT INTO organizations (id, name, owner_user_id, created_at) VALUES (?, ?, ?, ?)`,
			orgID, req.OrgName, userID, now); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		if _, err := tx.Exec(`INSERT INTO users (id, email, password_hash, org_id, created_at) VALUES (?, ?, ?, ?, ?)`,
			userID, req.Email, hash, orgID, now); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		if err := tx.Commit(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}

		token, err := auth.GenerateToken(cfg.JWTSecret, userID, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "token error")
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":    userID,
				"email": req.Email,
			},
			"org": fiber.Map{
				"id":   orgID,
				"name": req.OrgName,
			},
		})
	}
}

func handleLogin(db *sql.DB, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req loginRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
		if req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing required fields")
		}

		var userID, hash, orgID string
		if err := db.QueryRow(`SELECT id, password_hash, org_id FROM users WHERE email = ?`, req.Email).Scan(&userID, &hash, &orgID); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		if err := auth.CheckPassword(hash, req.Password); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		token, err := auth.GenerateToken(cfg.JWTSecret, userID, orgID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "token error")
		}
		return c.JSON(fiber.Map{"token": token})
	}
}

func handleMe(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := userIDFrom(c)
		orgID := orgIDFrom(c)
		var email string
		if err := db.QueryRow(`SELECT email FROM users WHERE id = ?`, userID).Scan(&email); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		var orgName string
		if err := db.QueryRow(`SELECT name FROM organizations WHERE id = ?`, orgID).Scan(&orgName); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "org not found")
		}
		return c.JSON(fiber.Map{
			"user": fiber.Map{"id": userID, "email": email},
			"org":  fiber.Map{"id": orgID, "name": orgName},
		})
	}
}
