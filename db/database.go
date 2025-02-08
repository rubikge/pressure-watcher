package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "pressure_db"
)

func InitDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS pressure_logs (
                id SERIAL PRIMARY KEY,
                timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                min_pressure REAL,
                max_pressure REAL,
                significant_change BOOLEAN
        );
        `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return db, nil
}

func SavePressureLog(db *sql.DB, minPressure, maxPressure float64, significantChange bool) {
	insertQuery := `
        INSERT INTO pressure_logs (min_pressure, max_pressure, significant_change)
        VALUES ($1, $2, $3);
        `
	_, err := db.Exec(insertQuery, minPressure, maxPressure, significantChange)
	if err != nil {
		log.Println("Error saving to database:", err)
	}
}

func GetLastPressureLog(db *sql.DB) (float64, float64, bool, error) {
	query := `SELECT min_pressure, max_pressure, significant_change FROM pressure_logs ORDER BY timestamp DESC LIMIT 1;`
	var minPressure float64
	var maxPressure float64
	var significantChange bool
	err := db.QueryRow(query).Scan(&minPressure, &maxPressure, &significantChange)
	return minPressure, maxPressure, significantChange, err
}
