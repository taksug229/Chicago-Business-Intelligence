package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelvins/geocoder"
	_ "github.com/lib/pq"

	"main/database"
	"main/models"
)

func InitializeAPIKey() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	geocoder.ApiKey = os.Getenv("API_KEY")
}

func fetchData(db *sql.DB, url string, limit string, target interface{}, insertFunc func(*sql.DB, interface{}), dataSource string, tableName string) {
	log.Printf("Getting %s\n", dataSource)
	var fullURL string
	if limit != "" {
		fullURL = fmt.Sprintf("%s?$limit=%s", url, limit)
	} else {
		fullURL = url
	}
	res, err := http.Get(fullURL)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, target)
	insertFunc(db, target)
	log.Printf("Success! %s to table: %s\n", dataSource, tableName)
}

func GetOldTaxiTrips(db *sql.DB, limit string) {
	InitializeAPIKey()
	var taxiTripsList models.TaxiTripsJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/wrvz-psew.json", limit, &taxiTripsList, database.InsertTaxiTripsWrapper, "Taxi Trips from 2013-2023", "taxi_trips")
}

func GetOldRideShareTrips(db *sql.DB, limit string) {
	InitializeAPIKey()
	var taxiTripsList models.TaxiTripsJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/m6dm-c72p.json", limit, &taxiTripsList, database.InsertTaxiTripsWrapper, "Rideshare Trips from 2018-2022", "taxi_trips")
}

func GetNewTaxiTrips(db *sql.DB, limit string) {
	InitializeAPIKey()
	var taxiTripsList models.TaxiTripsJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/ajtu-isnz.json", limit, &taxiTripsList, database.InsertTaxiTripsWrapper, "New Taxi Trips", "taxi_trips")
}

func GetNewRideShareTrips(db *sql.DB, limit string) {
	InitializeAPIKey()
	var taxiTripsList models.TaxiTripsJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/n26f-ihde.json", limit, &taxiTripsList, database.InsertTaxiTripsWrapper, "New Rideshare Trips", "taxi_trips")
}

func GetCovidDaily(db *sql.DB, limit string) {
	var covidCasesList models.CovidDailyJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/naz8-j4nc.json", limit, &covidCasesList, database.InsertCovidDailyWrapper, "Covid Daily", "covid_daily")
}

func GetCovidWeeklyZipCode(db *sql.DB, limit string) {
	var covidCasesList models.CovidWeeklyZipJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/yhhz-zm2v.json", limit, &covidCasesList, database.InsertCovidWeeklyZipWrapper, "Covid Weekly", "covid_weekly_zip")
}

func GetCCVI(db *sql.DB, limit string) {
	var ccviCasesList models.CCVIJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/xhc6-88s9.json", limit, &ccviCasesList, database.InsertCCVIWrapper, "CCVI", "ccvi_ca & ccvi_zip")
}

func GetBuildingPermits(db *sql.DB, limit string) {
	var permitList models.BuildingPermitsJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/ydr8-5enu.json", limit, &permitList, database.InsertBuildingPermitsWrapper, "Building Permits", "building_permits")
}

func GetUnemploymentRates(db *sql.DB, limit string) {
	var ratesList models.UnemploymentRatesJsonRecords
	fetchData(db, "https://data.cityofchicago.org/resource/iqnk-2tcu.json", limit, &ratesList, database.InsertUnemploymentRatesWrapper, "Unemployment", "unemployment")
}
