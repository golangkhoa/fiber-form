package main

import (
	// "fmt"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/template/html/v2"
)

func main() {
	// engine := html.New("./views", ".tmpl")
    // app := fiber.New(fiber.Config{
	// 	Views: engine,
	// })
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
    })

	app.Listen(":3000")
}