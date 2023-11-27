package main

import (
	"log"

	"github.com/emarifer/gofiber-htmx-sessions/db"
	"github.com/emarifer/gofiber-htmx-sessions/routes"
	"github.com/emarifer/gofiber-htmx-sessions/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func init() {
	// Create/migrate DB
	db.MakeMigrations()
	// Init sessions store
	sessions.InitSessionsStore()
}

func main() {
	// We make sure that when leaving the connection with the DB is closed
	defer db.Db.Close()

	// Create a new engine
	engine := html.New("./views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "base.layout",
	})

	app.Static("/", "./assets")

	app.Use(logger.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

/* CENTRADO CON POSITION RESPONSIVE. VER:
https://stackoverflow.com/questions/1776915/how-can-i-center-an-absolutely-positioned-element-in-a-div#23384995

FLASH MESSAGES IN GO FIBER. VER:
https://github.com/sujit-baniya/flash

USO DE VARIABLES EN CICLOS "RANGE". VER:
https://stackoverflow.com/questions/43263280/go-template-cant-evaluate-field-x-in-type-y-x-not-part-of-y-but-stuck-in-a#43263399

REFERENCES:
https://blog.bytebytego.com/p/password-session-cookie-token-jwt
https://docs.gofiber.io/api/middleware/session
https://github.com/gofiber/storage/
https://github.com/gofiber/storage/tree/main/sqlite3
https://github.com/gofiber/recipes/tree/master/sessions-sqlite3

https://www.techtarget.com/searchdatamanagement/definition/RDBMS-relational-database-management-system

https://www.epochconverter.com/

https://stackoverflow.com/questions/1409649/how-to-change-the-height-of-a-br
https://stackoverflow.com/questions/1776915/how-can-i-center-an-absolutely-positioned-element-in-a-div#23384995
*/
