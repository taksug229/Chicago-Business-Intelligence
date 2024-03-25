package database

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"main/models"
	"main/utils"
)

func ClearDatabase(createNewDatabase bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	initial_db_connection := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable port=5432", user, password, host)
	db, err := sql.Open("postgres", initial_db_connection)
	if err != nil {
		log.Fatal(fmt.Println("Couldn't Open Connection to database"))
		panic(err)
	}
	defer db.Close()

	drop_database_query := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbname)
	_, err = db.Exec(drop_database_query)
	if err != nil {
		log.Fatal(err)
	}
	drop_message := fmt.Sprintf("Dropped %s database (if existed).", dbname)
	log.Println(drop_message)
	if createNewDatabase {
		create_database_query := fmt.Sprintf("CREATE DATABASE %s;", dbname)
		_, err = db.Exec(create_database_query)
		if err != nil {
			log.Fatal(err)
		}
		create_message := fmt.Sprintf("Created new  %s database.", dbname)
		log.Println(create_message)
	}
}

func Login2Database() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	dbConnection := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable port=5432", user, dbname, password, host)
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal("Couldn't Open Connection to database")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("Couldn't Connect to database")
		panic(err)
	}
	log.Println("Connected to database.")
	return db
}

func RunQueriesFromFile(db *sql.DB, filename string) error {
	queryBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	queries := strings.Split(string(queryBytes), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query != "" {
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
	}
	message := fmt.Sprintf("Success! %s", filename)
	log.Println(message)
	return nil
}

const (
	maxRetries     = 5
	initialBackoff = 10 * time.Second
)

func isDatabaseContendedError(err error) bool {
	return strings.Contains(err.Error(), "pq: database") && strings.Contains(err.Error(), "is being accessed by other users")
}

func insertOrUpdateDataWithRetry(db *sql.DB, tableName string, data map[string]interface{}, conflictColumns []string) error {
	backoff := initialBackoff
	for i := 0; i < maxRetries; i++ {
		err := insertOrUpdateData(db, tableName, data, conflictColumns)
		if err == nil {
			return nil
		}

		// If the error is not related to database contention, return immediately
		if !isDatabaseContendedError(err) {
			return err
		}

		// Exponential backoff before retrying
		time.Sleep(backoff)
		backoff *= 2
	}
	return errors.New("max retries exceeded, unable to insert/update data")
}

func insertOrUpdateData(db *sql.DB, tableName string, data map[string]interface{}, conflictColumns []string) error {
	columns := make([]string, 0)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)

	var i int = 1
	for column, value := range data {
		columns = append(columns, column)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		values = append(values, value)
		i++
	}

	insertQuery := "INSERT INTO " + tableName + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ") ON CONFLICT (" + strings.Join(conflictColumns, ", ") + ") DO UPDATE SET "
	updateValues := make([]string, 0)
	for _, column := range columns {
		if !contains(conflictColumns, column) {
			updateValues = append(updateValues, column+"=EXCLUDED."+column)
		}
	}
	insertQuery += strings.Join(updateValues, ", ")
	_, err := db.Exec(insertQuery, values...)
	return err
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func InsertZip2CAFromCSV(db *sql.DB, csvfile string) error {
	file, err := os.Open(csvfile)
	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Println("Error reading CSV:", err)
		return err
	}
	for _, record := range records {
		if record[0] == "zip_code" {
			continue
		}
		caMap := make(map[string]sql.NullString)
		for i := 1; i < 10; i++ {
			key := fmt.Sprintf("ca%d", i)
			if utils.StringMissing(record[i]) {
				caMap[key] = sql.NullString{}
			} else {
				caMap[key] = sql.NullString{String: record[i], Valid: true}
			}
		}

		data := map[string]interface{}{
			"zip_code":        record[0],
			"community_area1": caMap["ca1"],
			"community_area2": caMap["ca2"],
			"community_area3": caMap["ca3"],
			"community_area4": caMap["ca4"],
			"community_area5": caMap["ca5"],
			"community_area6": caMap["ca6"],
			"community_area7": caMap["ca7"],
			"community_area8": caMap["ca8"],
			"community_area9": caMap["ca9"],
		}

		conflictColumns := []string{"zip_code"}
		err := insertOrUpdateData(db, "zip2ca", data, conflictColumns)
		if err != nil {
			log.Println("Error inserting data:", err)
			return err
		}
	}
	log.Println("Success! Table: zip2ca")
	return nil
}

func InsertTaxiTrips(db *sql.DB, trips_list models.TaxiTripsJsonRecords) {
	for i := 0; i < len(trips_list); i++ {
		trip_id := trips_list[i].Trip_id
		if utils.StringMissing(trip_id) {
			continue
		}
		taxi_id := sql.NullString{}
		trip_type := "rideshare"
		if !utils.StringMissing(trips_list[i].Taxi_id) {
			taxi_id = sql.NullString{String: trips_list[i].Taxi_id, Valid: true}
			trip_type = "taxi"
		}

		trip_start_timestamp := trips_list[i].Trip_start_timestamp
		if len(trip_start_timestamp) < 23 {
			continue
		}
		trip_end_timestamp := trips_list[i].Trip_end_timestamp
		if len(trip_end_timestamp) < 23 {
			continue
		}
		pickup_centroid_latitude := trips_list[i].Pickup_centroid_latitude
		if utils.StringMissing(pickup_centroid_latitude) {
			continue
		}
		pickup_centroid_longitude := trips_list[i].Pickup_centroid_longitude
		if utils.StringMissing(pickup_centroid_longitude) {
			continue
		}
		pickup_community_area := trips_list[i].Pickup_community_area
		if utils.StringMissing(pickup_community_area) {
			continue
		}
		dropoff_centroid_latitude := trips_list[i].Dropoff_centroid_latitude
		if utils.StringMissing(dropoff_centroid_latitude) {
			continue
		}
		dropoff_centroid_longitude := trips_list[i].Dropoff_centroid_longitude
		if utils.StringMissing(dropoff_centroid_longitude) {
			continue
		}
		dropoff_community_area := trips_list[i].Dropoff_community_area
		if utils.StringMissing(dropoff_community_area) {
			continue
		}
		pickup_lat, _ := strconv.ParseFloat(pickup_centroid_latitude, 64)
		pickup_lng, _ := strconv.ParseFloat(pickup_centroid_longitude, 64)
		dropoff_lat, _ := strconv.ParseFloat(dropoff_centroid_latitude, 64)
		dropoff_lng, _ := strconv.ParseFloat(dropoff_centroid_longitude, 64)
		pickup_zip := utils.ReverseGeocodeToZipCode(pickup_lat, pickup_lng)
		if utils.StringMissing(pickup_zip) {
			continue
		}
		dropoff_zip := utils.ReverseGeocodeToZipCode(dropoff_lat, dropoff_lng)
		if utils.StringMissing(dropoff_zip) {
			continue
		}
		data := map[string]interface{}{
			"trip_id":                    trip_id,
			"trip_type":                  trip_type,
			"taxi_id":                    taxi_id,
			"trip_start_timestamp":       trip_start_timestamp,
			"trip_end_timestamp":         trip_end_timestamp,
			"pickup_centroid_latitude":   pickup_centroid_latitude,
			"pickup_centroid_longitude":  pickup_centroid_longitude,
			"pickup_community_area":      pickup_community_area,
			"dropoff_centroid_latitude":  dropoff_centroid_latitude,
			"dropoff_centroid_longitude": dropoff_centroid_longitude,
			"dropoff_community_area":     dropoff_community_area,
			"pickup_zip_code":            pickup_zip,
			"dropoff_zip_code":           dropoff_zip,
		}

		conflictColumns := []string{"trip_id", "trip_type"}
		err := insertOrUpdateDataWithRetry(db, "taxi_trips", data, conflictColumns)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func InsertCovidDaily(db *sql.DB, covid_list models.CovidDailyJsonRecords) {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i := 0; i < len(covid_list); i++ {
		currentTime := time.Now().In(location)
		update_date := currentTime.Format("2006-01-02T15:04:05.000")

		report_date := covid_list[i].Date
		if utils.StringMissing(report_date) {
			continue
		}
		cases := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Cases) {
			cases = sql.NullString{String: covid_list[i].Cases, Valid: true}
		}
		deaths := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Deaths) {
			deaths = sql.NullString{String: covid_list[i].Deaths, Valid: true}
		}
		hospitalizations := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Hospitalizations) {
			hospitalizations = sql.NullString{String: covid_list[i].Hospitalizations, Valid: true}
		}
		data := map[string]interface{}{
			"report_date":      report_date,
			"update_date":      update_date,
			"cases":            cases,
			"deaths":           deaths,
			"hospitalizations": hospitalizations,
		}

		conflictColumns := []string{"report_date"}
		err := insertOrUpdateData(db, "covid_daily", data, conflictColumns)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func InsertCovidWeeklyZip(db *sql.DB, covid_list models.CovidWeeklyZipJsonRecords) {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i := 0; i < len(covid_list); i++ {
		currentTime := time.Now().In(location)
		update_date := currentTime.Format("2006-01-02T15:04:05.000")
		zip_code := covid_list[i].Zip_code
		if utils.StringMissing(zip_code) {
			continue
		}
		week_number := covid_list[i].Week_number
		if utils.StringMissing(week_number) {
			continue
		}
		week_start := covid_list[i].Week_start
		if utils.StringMissing(week_start) {
			continue
		}
		week_end := covid_list[i].Week_end
		if utils.StringMissing(week_end) {
			continue
		}
		cases := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Cases) {
			cases = sql.NullString{String: covid_list[i].Cases, Valid: true}
		}
		cases_rate := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Cases_rate) {
			cases_rate = sql.NullString{String: covid_list[i].Cases_rate, Valid: true}
		}
		tests := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Tests) {
			tests = sql.NullString{String: covid_list[i].Tests, Valid: true}
		}
		tests_rate := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Tests_rate) {
			tests_rate = sql.NullString{String: covid_list[i].Tests_rate, Valid: true}
		}
		percent_tested_positive := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Percent_tested_positive) {
			percent_tested_positive = sql.NullString{String: covid_list[i].Percent_tested_positive, Valid: true}
		}
		deaths := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Deaths) {
			deaths = sql.NullString{String: covid_list[i].Deaths, Valid: true}
		}
		deaths_rate := sql.NullString{}
		if !utils.StringMissing(covid_list[i].Deaths_rate) {
			deaths_rate = sql.NullString{String: covid_list[i].Deaths_rate, Valid: true}
		}
		data := map[string]interface{}{
			"zip_code":                zip_code,
			"week_number":             week_number,
			"week_start":              week_start,
			"week_end":                week_end,
			"update_date":             update_date,
			"cases":                   cases,
			"cases_rate":              cases_rate,
			"tests":                   tests,
			"tests_rate":              tests_rate,
			"percent_tested_positive": percent_tested_positive,
			"deaths":                  deaths,
			"deaths_rate":             deaths_rate,
		}

		conflictColumns := []string{"zip_code", "week_start", "week_end"}
		err := insertOrUpdateData(db, "covid_weekly_zip", data, conflictColumns)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func InsertCCVI(db *sql.DB, ccvi_list models.CCVIJsonRecords) {
	for i := 0; i < len(ccvi_list); i++ {
		geography_type := ccvi_list[i].Geography_type
		ccvi_score := ccvi_list[i].CCVI_score
		ccvi_category := ccvi_list[i].CCVI_category
		if utils.StringMissing(ccvi_category) {
			continue
		}
		if geography_type == "CA" {
			community_area := ccvi_list[i].Community_or_zip
			community_area_name := ccvi_list[i].Community_area_name

			data := map[string]interface{}{
				"community_area":      community_area,
				"community_area_name": community_area_name,
				"ccvi_score":          ccvi_score,
				"ccvi_category":       ccvi_category,
			}

			conflictColumns := []string{"community_area"}
			err := insertOrUpdateData(db, "ccvi_ca", data, conflictColumns)
			if err != nil {
				log.Print(err)
				continue
			}
		} else if geography_type == "ZIP" {
			zip_code := ccvi_list[i].Community_or_zip
			data := map[string]interface{}{
				"zip_code":      zip_code,
				"ccvi_score":    ccvi_score,
				"ccvi_category": ccvi_category,
			}

			conflictColumns := []string{"zip_code"}
			err := insertOrUpdateData(db, "ccvi_zip", data, conflictColumns)
			if err != nil {
				log.Print(err)
				continue
			}
		} else {
			continue
		}
	}
}

func InsertBuildingPermits(db *sql.DB, building_list models.BuildingPermitsJsonRecords) {
	for i := 0; i < len(building_list); i++ {
		id := building_list[i].Id
		if utils.StringMissing(id) {
			continue
		}
		permit_nbr := building_list[i].Permit_nbr
		if utils.StringMissing(permit_nbr) {
			continue
		}
		zip_code := building_list[i].Zip_code
		if utils.StringMissing(zip_code) {
			continue
		}
		zip_code = strings.TrimPrefix(zip_code, "-")
		zip_code = strings.TrimSuffix(zip_code, "-")
		zip_code = strings.Replace(zip_code, "-", "", -1)
		zip_code = zip_code[:5]
		permit_type := building_list[i].Permit_type
		if utils.StringMissing(permit_type) {
			continue
		}
		review_type := building_list[i].Review_type
		if utils.StringMissing(review_type) {
			continue
		}
		application_date := building_list[i].Application_date
		if utils.StringMissing(application_date) {
			continue
		}
		issue_date := building_list[i].Issue_date
		if utils.StringMissing(issue_date) {
			continue
		}
		community_area := sql.NullString{}
		if !utils.StringMissing(building_list[i].Community_area) {
			community_area = sql.NullString{String: building_list[i].Community_area, Valid: true}
		}
		total_fee := sql.NullString{}
		if !utils.StringMissing(building_list[i].Total_fee) {
			total_fee = sql.NullString{String: building_list[i].Total_fee, Valid: true}
		}
		data := map[string]interface{}{
			"id":               id,
			"permit_nbr":       permit_nbr,
			"permit_type":      permit_type,
			"review_type":      review_type,
			"application_date": application_date,
			"issue_date":       issue_date,
			"total_fee":        total_fee,
			"zip_code":         zip_code,
			"community_area":   community_area,
		}

		conflictColumns := []string{"id", "permit_nbr"}
		err := insertOrUpdateData(db, "building_permits", data, conflictColumns)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
func InsertUnemploymentRates(db *sql.DB, unemployment_list models.UnemploymentRatesJsonRecords) {
	for i := 0; i < len(unemployment_list); i++ {
		community_area := unemployment_list[i].Community_area
		if utils.StringMissing(community_area) {
			continue
		}
		community_area_name := unemployment_list[i].Community_area_name
		if utils.StringMissing(community_area_name) {
			continue
		}
		unemployment := unemployment_list[i].Unemployment
		if utils.StringMissing(unemployment) {
			continue
		}
		poverty := sql.NullString{}
		if !utils.StringMissing(unemployment_list[i].Poverty) {
			poverty = sql.NullString{String: unemployment_list[i].Poverty, Valid: true}
		}
		income := sql.NullString{}
		if !utils.StringMissing(unemployment_list[i].Income) {
			income = sql.NullString{String: unemployment_list[i].Income, Valid: true}
		}
		data := map[string]interface{}{
			"community_area":      community_area,
			"community_area_name": community_area_name,
			"unemployment_rate":   unemployment,
			"poverty_rate":        poverty,
			"per_capita_income":   income,
		}

		conflictColumns := []string{"community_area"}
		err := insertOrUpdateData(db, "unemployment", data, conflictColumns)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
