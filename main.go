package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"	
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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
	readDB, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("signup", fiber.Map{})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		insert, _ := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
		_, err = insert.Exec(user.Username, user.Password)
		if err != nil {
			fmt.Println(err)
		}

		return c.Render("signup", fiber.Map{})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		readDB.First(&user, "username = ?", user.Username)
		readDB.First(&user, "password = ?", user.Password)

		return c.Render("login", fiber.Map{})
	})
	
	app.Listen(":3000")
	db.Close()
}
