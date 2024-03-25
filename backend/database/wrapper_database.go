package database

import (
	"database/sql"
	"main/models"
)

func InsertTaxiTripsWrapper(db *sql.DB, data interface{}) {
	tripsList := data.(*models.TaxiTripsJsonRecords)
	InsertTaxiTrips(db, *tripsList)
}

func InsertCovidDailyWrapper(db *sql.DB, data interface{}) {
	covidList := *data.(*models.CovidDailyJsonRecords)
	InsertCovidDaily(db, covidList)
}

func InsertCovidWeeklyZipWrapper(db *sql.DB, data interface{}) {
	covidList := *data.(*models.CovidWeeklyZipJsonRecords)
	InsertCovidWeeklyZip(db, covidList)
}

func InsertCCVIWrapper(db *sql.DB, data interface{}) {
	ccviList := *data.(*models.CCVIJsonRecords)
	InsertCCVI(db, ccviList)
}

func InsertBuildingPermitsWrapper(db *sql.DB, data interface{}) {
	permitList := *data.(*models.BuildingPermitsJsonRecords)
	InsertBuildingPermits(db, permitList)
}

func InsertUnemploymentRatesWrapper(db *sql.DB, data interface{}) {
	ratesList := *data.(*models.UnemploymentRatesJsonRecords)
	InsertUnemploymentRates(db, ratesList)
}
