package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
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

// Checks if the user exists in the database
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

// Does other things such as: Loading template files,
// encrypting username and password, etc.
func main() {
	var user User
	engine := html.New("./views", ".tmpl")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	Encrptedusername := sha512.Sum512([]byte(user.Username))
	Encryptedpassword := sha512.Sum512([]byte(user.Password))
	fmt.Println(Encrptedusername)
	fmt.Println(Encryptedpassword)

	db, err := sql.Open("sqlite3", "user.db")
	if err != nil {
		panic("failed to open database")
	}
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

		if !UserExists(db, hex.EncodeToString(Encrptedusername[:])) {
			insert, _ := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
			_, err = insert.Exec(hex.EncodeToString(Encrptedusername[:]), hex.EncodeToString(Encryptedpassword[:]))
			if err != nil {
				fmt.Println(err)
			}

			return c.Render("signup", fiber.Map{
				"Success": "Signup Succussful",
			})
		}

		return c.Render("signup", fiber.Map{})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		if err := c.BodyParser(&user); err != nil {
			return err
		}

		if UserExists(db, hex.EncodeToString(Encrptedusername[:])) {
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
