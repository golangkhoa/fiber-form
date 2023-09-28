package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string `form:"name"`
	Password string `form:"pwd"`
}

func main() {
	var user User
	engine := html.New("./views", ".tmpl")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	db, err := sql.Open("sqlite3", "user.db")
	if err != nil {
		panic("failed to connect database")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		insert, _ := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
		_, err = insert.Exec(user.Username, user.Password)
		fmt.Println(err)

		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Listen(":3000")
	db.Close()
}
