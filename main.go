package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
	"os"
)
// useless comment
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

	_, err := db.Exec("INSERT into test VALUES (toto)")
	if err != nil {
		log.Fatalf("An error occured while executing query: %v", err)
	}

	return c.Redirect("/")
}

func main() {
	log.Print("TOTO DEMARRE")
	log.Print(os.Getenv("QOVERY_POSTGRESQL_ZF0F42794_DATABASE_URL_INTERNAL") + "?sslmode=disable")
	db, err := sql.Open("postgres", os.Getenv("QOVERY_POSTGRESQL_DATABASE_URL_INTERNAL")+"?sslmode=disable")
	if err != nil {
		log.Print(os.Getenv("QOVERY_POSTGRESQL_ZF0F42794_DATABASE_URL_INTERNAL"))
		log.Print(err)
	}

	_, createErr := db.Exec("CREATE TABLE IF NOT EXISTS test (id serial PRIMARY KEY, username VARCHAR ( 50 ) NOT NULL);")
	if createErr != nil {
		log.Fatal(createErr)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		log.Printf("Posting :%s")
		return postHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
