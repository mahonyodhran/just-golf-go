package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	logger *log.Logger
)

const (
	dbConnEnv = "DB_CONN"
	driver    = "postgres"
)

const (
	scorecardTable = "scorecard"
	courseTable    = "course"
	golferTable    = "golfer"
)

const (
	queryGetScorecards        = `SELECT id, hole_1, hole_2, hole_3, hole_4, hole_5, hole_6, hole_7, hole_8, hole_9, hole_10, hole_11, hole_12, hole_13, hole_14, hole_15, hole_16, hole_17, hole_18 FROM ` + scorecardTable
	queryInsertScore          = `INSERT INTO ` + scorecardTable + ` (date_inserted, golfer_id, course_id, hole_1, hole_2, hole_3, hole_4, hole_5, hole_6, hole_7, hole_8, hole_9, hole_10, hole_11, hole_12, hole_13, hole_14, hole_15, hole_16, hole_17, hole_18) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`
	queryCreateScorecardTable = `CREATE TABLE IF NOT EXISTS ` + scorecardTable + ` (id SERIAL PRIMARY KEY, date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP, golfer_id INT REFERENCES ` + golferTable + `(id), course_id INT REFERENCES ` + courseTable + `(id), hole_1 INT, hole_2 INT, hole_3 INT, hole_4 INT, hole_5 INT, hole_6 INT, hole_7 INT, hole_8 INT, hole_9 INT, hole_10 INT, hole_11 INT, hole_12 INT, hole_13 INT, hole_14 INT, hole_15 INT, hole_16 INT, hole_17 INT, hole_18 INT)`

	queryGetCourses        = `SELECT ID, NAME FROM ` + courseTable
	queryInsertCourse      = `INSERT INTO ` + courseTable + ` (name) VALUES ($1)`
	queryCreateCourseTable = `CREATE TABLE IF NOT EXISTS ` + courseTable + ` (id SERIAL PRIMARY KEY, date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP, name TEXT)`

	queryCreateGolferTable = `CREATE TABLE IF NOT EXISTS ` + golferTable + ` (id SERIAL PRIMARY KEY, date_inserted TIMESTAMP DEFAULT CURRENT_TIMESTAMP, first_name TEXT, last_name TEXT, index int)`
)

type Scorecard struct {
	ID       int
	GolferID int
	CourseID int
	Holes    [18]int
}

type Course struct {
	ID   int
	Name string
}

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	connStr := os.Getenv(dbConnEnv)

	db, err = sql.Open(driver, connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// err = dropTables()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = createCourseTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createGolferTable()
	if err != nil {
		log.Fatal(err)
	}

	err = createScorecardTable()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	logger = log.New(logFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getScorecards() ([]Scorecard, error) {
	data := []Scorecard{}
	rows, err := db.Query(queryGetScorecards)
	if err != nil {
		return nil, fmt.Errorf("error running query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var scorecard Scorecard
		err := rows.Scan(&scorecard.ID, &scorecard.Holes[0], &scorecard.Holes[1], &scorecard.Holes[2], &scorecard.Holes[3], &scorecard.Holes[4], &scorecard.Holes[5], &scorecard.Holes[6], &scorecard.Holes[7], &scorecard.Holes[8], &scorecard.Holes[9], &scorecard.Holes[10], &scorecard.Holes[11], &scorecard.Holes[12], &scorecard.Holes[13], &scorecard.Holes[14], &scorecard.Holes[15], &scorecard.Holes[16], &scorecard.Holes[17])
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, scorecard)
	}
	return data, nil
}

func insertScore(scorecard Scorecard) error {
	_, err := db.Exec(queryInsertScore, time.Now(), scorecard.GolferID, scorecard.CourseID, scorecard.Holes[0], scorecard.Holes[1], scorecard.Holes[2], scorecard.Holes[3],
		scorecard.Holes[4], scorecard.Holes[5], scorecard.Holes[6], scorecard.Holes[7], scorecard.Holes[8],
		scorecard.Holes[9], scorecard.Holes[10], scorecard.Holes[11], scorecard.Holes[12], scorecard.Holes[13],
		scorecard.Holes[14], scorecard.Holes[15], scorecard.Holes[16], scorecard.Holes[17])
	if err != nil {
		logger.Println("Error inserting scorecard:", err)
	} else {
		logger.Printf("Inserted record to Scorecard (ID: %d)", scorecard.ID)
	}
	return err
}

func createScorecardTable() error {
	_, err := db.Exec(queryCreateScorecardTable)
	if err != nil {
		logger.Println("Error creating scorecard table", err)
	}
	return err
}

func createCourseTable() error {
	_, err := db.Exec(queryCreateCourseTable)
	if err != nil {
		logger.Println("Error creating course table", err)
	}
	return err
}

func createGolferTable() error {
	_, err := db.Exec(queryCreateGolferTable)
	if err != nil {
		logger.Printf("Error creating %s table: %v", golferTable, err)
	}
	return err
}

func getCourses() ([]Course, error) {
	data := []Course{}
	rows, err := db.Query(queryGetCourses)
	if err != nil {
		log.Fatal("Error running query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Name)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, course)
	}
	return data, nil
}

func insertCourse(courseName string) error {
	_, err := db.Exec(queryInsertCourse, courseName)
	if err != nil {
		logger.Println("Error inserting course:", err)
	} else {
		logger.Printf("Inserted record to Course (Name: %s)", courseName)
	}
	return err
}

// func dropTables() error {
// 	//Obviously use with caution - take an extract beforehand
// 	_, err := db.Exec(`DROP TABLE if exists course, golfer, scorecard,`)
// 	if err != nil {
// 		logger.Println("Error dropping tables", err)
// 	}
// 	return err
// }
