package model

import (
	"fmt"
	"log"
	"os"
	
	"encoding/json"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
)

func dbConn() *pgx.ConnPool {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	pgxConf, err := pgx.ParseConnectionString(connString)

	if err != nil {
		log.Println(err)
		return nil
	}

	pgxPollConf := pgx.ConnPoolConfig{
		ConnConfig:     pgxConf,
		MaxConnections: 5,
	}

	db, err := pgx.NewConnPool(pgxPollConf)
	if err != nil {
		log.Println(err)
		return nil
	}
  
  return db
}

func getMovie(c *fiber.Ctx) error {
  	db := dbConn() // inisialisasi koneksi db yang berulang
  	defer db.Close()
  
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var movies []*Movie

	for rows.Next() {
		var movie *Movie
		var updatedAt time.Time
		err := rows.Scan(
			&movie.ID,
			&movie.UserID,
			&movie.Title,
			&movie.Year,
			&movie.Rating,
			&movie.Trailer,
			&movie.Synopsis,
			&movie.Info,
			&movie.ImageLink,
			&movie.CoverLink,
			&movie.Views,
			&movie.Active,
			&movie.Ongoing,
			&movie.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		movies = append(movies, movie)
	}

	return c.JSON(movies)
}

func main() {
  	app := fiber.New()
  	app.get("/movie", getMovie)
	app.Listen(":3000")
}
