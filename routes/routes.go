package routes

import (
	"github.com/emarifer/gofiber-htmx-sessions/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	/* Views */
	app.Get("/", handlers.RedirectToLogin)
	app.Get("/login", handlers.HandleViewLogin)
	app.Get("/register", handlers.HandleViewRegister)
	app.Get("/profile", handlers.HandleViewProfile)

	/* Singin/Singup & Sessions  */
	app.Post("/api/signin", handlers.HandleSigninUser)
	app.Post("/api/signup", handlers.HandleRegisterUser)
	app.Post("/api/logout", handlers.HandleSessionLogout)
}
