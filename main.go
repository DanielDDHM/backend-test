package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/database_actions"
	"github.com/japhy-tech/backend-test/internal"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
	CSVFile  = "breeds.csv"
)

func main() {
	logger := charmLog.NewWithOptions(os.Stderr, charmLog.Options{
		Formatter:       charmLog.TextFormatter,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "🧑‍💻 backend-test",
		Level:           charmLog.DebugLevel,
	})

	err := database_actions.InitMigrator(MysqlDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}

	msg, err := database_actions.RunMigrate("up", 0)
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info(msg)
	}

	db, err := sql.Open("mysql", MysqlDSN)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxIdleConns(0)

	err = db.Ping()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	logger.Info("Database connected")

	err = ensureDataInDatabase(db, logger)

	app := internal.NewApp(logger)

	r := mux.NewRouter()
	app.RegisterRoutes(r.PathPrefix("/v1").Subrouter())

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	err = http.ListenAndServe(
		net.JoinHostPort("", ApiPort),
		r,
	)

	// =============================== Starting Msg ===============================
	logger.Info(fmt.Sprintf("Service started and listen on port %s", ApiPort))
}

func ensureDataInDatabase(db *sql.DB, logger *charmLog.Logger) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM breeds").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking if 'breeds' table is empty: %s", err.Error())
	}

	if count == 0 {

		file, err := os.Open(CSVFile)
		if err != nil {
			return fmt.Errorf("error opening CSV file '%s': %s", CSVFile, err.Error())
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = ','

		lines, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("error reading CSV content: %s", err.Error())
		}

		stmt, err := db.Prepare("INSERT INTO breeds (id, species, pet_size, pet_name, average_male_adult_weight, average_female_adult_weight) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			return fmt.Errorf("error preparing insert statement: %s", err.Error())
		}
		defer stmt.Close()

		for _, line := range lines[1:] {
            id, _ := strconv.Atoi(strings.Trim(line[0], `"`))
            species := strings.Trim(line[1], `"`)
            pet_size := strings.Trim(line[2], `"`)
            pet_name := strings.Trim(line[3], `"`)
            average_male_adult_weight, _ := strconv.Atoi(strings.Trim(line[4], `"`))
            average_female_adult_weight, _ := strconv.Atoi(strings.Trim(line[5], `"`))

			_, err := stmt.Exec(id, species, pet_size, pet_name, average_male_adult_weight, average_female_adult_weight)
            if err != nil {
                logger.Error(fmt.Sprintf("Error inserting data for line %s: %s", line, err.Error()))
                continue
            }
		}
	}

	return nil
}