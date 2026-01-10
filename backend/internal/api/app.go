package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"nudgepay/internal/config"
)

func NewApp(db *sql.DB, cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: jsonErrorHandler,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/api/auth/register", handleRegister(db, cfg))
	app.Post("/api/auth/login", handleLogin(db, cfg))

	secured := app.Group("/api", authRequired(cfg.JWTSecret))
	secured.Get("/me", handleMe(db))
	secured.Get("/org", handleGetOrg(db))
	secured.Put("/org", handleUpdateOrg(db))

	secured.Get("/metrics", handleMetrics(db))

	secured.Get("/clients", handleListClients(db))
	secured.Post("/clients", handleCreateClient(db))
	secured.Get("/clients/:id", handleGetClient(db))
	secured.Put("/clients/:id", handleUpdateClient(db))
	secured.Delete("/clients/:id", handleDeleteClient(db))

	secured.Get("/templates", handleListTemplates(db))
	secured.Post("/templates", handleCreateTemplate(db))
	secured.Put("/templates/:id", handleUpdateTemplate(db))
	secured.Delete("/templates/:id", handleDeleteTemplate(db))

	secured.Get("/invoices", handleListInvoices(db))
	secured.Post("/invoices", handleCreateInvoice(db))
	secured.Get("/invoices/:id", handleGetInvoice(db))
	secured.Put("/invoices/:id", handleUpdateInvoice(db))
	secured.Delete("/invoices/:id", handleDeleteInvoice(db))

	secured.Get("/reminders", handleListReminders(db))
	secured.Post("/reminders/:id/send", handleSendReminder(db))
	secured.Post("/reminders/send-due", handleSendDueReminders(db))

	secured.Get("/outbox", handleListOutbox(db))

	return app
}
