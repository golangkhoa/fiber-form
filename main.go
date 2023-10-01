package main

import (
	"database/sql"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type User struct {
	Username string `form:"name"`
	Password string `form:"pwd"`
}

func UserExists(db *sql.DB, username string) bool {
	sqlStmt := `SELECT username FROM users WHERE username = ?`
	var u string
	err := db.QueryRow(sqlStmt, username).Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
	}
	return true
}

func main() {
	var user User
	engine := html.New("./views", ".tmpl")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	db, err := sql.Open("sqlite3", "user.db")
	defer db.Close()
	readDB, errDB := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
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

		if !UserExists(db, user.Username) {
			insert, _ := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
			_, err = insert.Exec(user.Username, user.Password)
			if err != nil {
				fmt.Println(err)
			}
		}

		return c.Render("signup", fiber.Map{
			"Success": "Signup Succussful",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		if UserExists(db, user.Username) {
			readDB.First(&user, "username = ?", user.Username)
			readDB.First(&user, "password = ?", user.Password)
			
			return c.Render("login", fiber.Map{
				"Success": "Login Succussful",
			})
		}

		return c.Render("login", fiber.Map{})
	})

	if errDB != nil {
		fmt.Println(errDB)
	}

	app.Listen(":3000")
}
