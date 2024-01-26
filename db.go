package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

var logger *log.Logger

type Scorecard struct {
	Holes [18]int
}

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	connStr := os.Getenv("DB_CONN")

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = createScorecardTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createCourseTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createGolferTable()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	logger = log.New(logFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func insertScore(scorecard Scorecard) error {
	_, err := db.Exec(`
		INSERT INTO scorecard (date_inserted, hole_1, hole_2, hole_3, hole_4, hole_5, hole_6, hole_7, hole_8, hole_9,
			hole_10, hole_11, hole_12, hole_13, hole_14, hole_15, hole_16, hole_17, hole_18)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`, time.Now(), scorecard.Holes[0], scorecard.Holes[1], scorecard.Holes[2], scorecard.Holes[3],
		scorecard.Holes[4], scorecard.Holes[5], scorecard.Holes[6], scorecard.Holes[7], scorecard.Holes[8],
		scorecard.Holes[9], scorecard.Holes[10], scorecard.Holes[11], scorecard.Holes[12], scorecard.Holes[13],
		scorecard.Holes[14], scorecard.Holes[15], scorecard.Holes[16], scorecard.Holes[17])

	if err != nil {
		logger.Println("Error inserting scorecard:", err)
	} else {
		logger.Println("Inserted record to Scorecard table")
	}

	return err
}

func createScorecardTable() error {
	_, err := db.Exec(`CREATE TABLE if not exists scorecard (
		id SERIAL PRIMARY KEY,
		date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		hole_1 INT,
		hole_2 INT,
		hole_3 INT,
		hole_4 INT,
		hole_5 INT,
		hole_6 INT,
		hole_7 INT,
		hole_8 INT,
		hole_9 INT,
		hole_10 INT,
		hole_11 INT,
		hole_12 INT,
		hole_13 INT,
		hole_14 INT,
		hole_15 INT,
		hole_16 INT,
		hole_17 INT,
		hole_18 INT
	)`)
	if err != nil {
		logger.Println("Error creating scorecard table", err)
	}

	return err
}

func createCourseTable() error {
	_, err := db.Exec(`CREATE TABLE if not exists course (
		id SERIAL PRIMARY KEY,
		date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT)`)
	if err != nil {
		logger.Println("Error creating course table", err)
	}
	return err
}

func createGolferTable() error {
	_, err := db.Exec(`CREATE TABLE if not exists golfer (
		id SERIAL PRIMARY KEY,
		date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		first_name TEXT,
		last_name TEXT,
		index int)`)
	if err != nil {
		logger.Println("Error creating golfer table", err)
	}
	return err
}