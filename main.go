package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var results []string
	rows, err := db.Query("SELECT * FROM test")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() {
		rows.Scan(&res)
		results = append(results, res)
	}
	return c.Render("index", fiber.Map{
		"Results": results,
	})
}
func postHandler(c *fiber.Ctx, db *sql.DB) error {
	input := ""
	if err := c.BodyParser(&input); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	fmt.Printf("%v", input)
	if input != "" {
		_, err := db.Exec("INSERT into test VALUES ($1)", input)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("PSQL"))
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}