package main

import (
	"database/sql"
	"log"
	"main/api"
	"main/database"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Get environment variables
	limit := os.Getenv("LIMIT")
	createTablesFilename := os.Getenv("CREATE_INITIAL_TABLE_FILE")
	zip2caCSVFile := os.Getenv("ZIP2CA_FILE")
	insertCATablesFilename := os.Getenv("INSERT_CANAME_TABLE_FILE")

	// Clear and login to the database
	database.ClearDatabase(true)
	db := database.Login2Database()
	defer db.Close()

	// Create Tables
	err = database.RunQueriesFromFile(db, createTablesFilename)
	if err != nil {
		log.Fatal(err)
	}
	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		err := database.InsertZip2CAFromCSV(db, zip2caCSVFile)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		api.GetCCVI(db, "")
	}()

	go func() {
		defer wg.Done()
		api.GetOldTaxiTrips(db, limit)
	}()

	go func() {
		defer wg.Done()
		api.GetOldRideShareTrips(db, limit)
	}()

	go func() {
		defer wg.Done()
		api.GetUnemploymentRates(db, "")
		err := database.RunQueriesFromFile(db, insertCATablesFilename)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	updateFreq := map[string]time.Duration{
		"GetNewTaxiTrips":       30 * 24 * time.Hour, // monthly
		"GetNewRideShareTrips":  30 * 24 * time.Hour, // monthly
		"GetCovidWeeklyZipCode": 7 * 24 * time.Hour,  // weekly
		"GetCovidDaily":         7 * 24 * time.Hour,  // weekly
		"GetBuildingPermits":    24 * time.Hour,      // daily
	}
	// Run the update loop indefinitely
	runUpdateLoop(db, limit, updateFreq)
}

func runUpdateLoop(db *sql.DB, limit string, updateFreq map[string]time.Duration) {
	for apiFunc, freq := range updateFreq {
		go func(apiFunc string, freq time.Duration) {
			for {
				// Call the respective API function
				switch apiFunc {
				case "GetNewTaxiTrips":
					api.GetNewTaxiTrips(db, limit)
				case "GetNewRideShareTrips":
					api.GetNewRideShareTrips(db, limit)
				case "GetCovidWeeklyZipCode":
					api.GetCovidWeeklyZipCode(db, limit)
				case "GetCovidDaily":
					api.GetCovidDaily(db, limit)
				case "GetBuildingPermits":
					api.GetBuildingPermits(db, limit)
				}
				// Determine port for HTTP service.
				port := os.Getenv("PORT")
				if port == "" {
					port = "8080"
					log.Printf("defaulting to port %s", port)
				}

				// Start HTTP server.
				log.Printf("listening on port %s", port)
				if err := http.ListenAndServe(":"+port, nil); err != nil {
					log.Println(err)
				}
				// Sleep for the specified duration before making the next call
				time.Sleep(freq)
			}
		}(apiFunc, freq)
	}
	select {}
}
