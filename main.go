package main

import (
	"fmt"
	"log"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `form:"name"`
	Password string `form:"pwd"`
}

func main() {
	var user User
	engine := html.New("./views", ".tmpl")
    app := fiber.New(fiber.Config{
		Views: engine,
	})
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		errDB := db.Where("username = ?", user.Username).Where("password = ?", user.Password).First(&user).Error

		if errDB != nil {
			if errDB == gorm.ErrRecordNotFound {
				err := db.Create(&User{Username: user.Username, Password: user.Password})

				if err != nil {
					fmt.Println(err)
				}
			} else {
				log.Fatalln(errDB)
			}
		}

		if err := c.BodyParser(&user); err != nil {
			return err
		}

		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	if user.ID > 1 && user.Username == user.Username {
		db.Delete(&user, 1)
	}
	if user.Username == "" && user.Password == "" {
		db.Delete(&user, 1)
	}

	app.Listen(":3000")
}