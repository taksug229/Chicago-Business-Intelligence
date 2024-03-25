package models

type TaxiTripsJsonRecords []struct {
	Trip_id                    string `json:"trip_id"`
	Taxi_id                    string `json:"taxi_id"`
	Trip_start_timestamp       string `json:"trip_start_timestamp"`
	Trip_end_timestamp         string `json:"trip_end_timestamp"`
	Pickup_centroid_latitude   string `json:"pickup_centroid_latitude"`
	Pickup_centroid_longitude  string `json:"pickup_centroid_longitude"`
	Pickup_community_area      string `json:"pickup_community_area"`
	Dropoff_centroid_latitude  string `json:"dropoff_centroid_latitude"`
	Dropoff_centroid_longitude string `json:"dropoff_centroid_longitude"`
	Dropoff_community_area     string `json:"dropoff_community_area"`
}

type CovidDailyJsonRecords []struct {
	Date             string `json:"lab_report_date"`
	Cases            string `json:"cases_total"`
	Deaths           string `json:"deaths_total"`
	Hospitalizations string `json:"hospitalizations_total"`
}

type CovidWeeklyZipJsonRecords []struct {
	Zip_code                string `json:"zip_code"`
	Week_number             string `json:"week_number"`
	Week_start              string `json:"week_start"`
	Week_end                string `json:"week_end"`
	Cases                   string `json:"cases_weekly"`
	Cases_rate              string `json:"case_rate_weekly"`
	Tests                   string `json:"tests_weekly"`
	Tests_rate              string `json:"test_rate_weekly"`
	Percent_tested_positive string `json:"percent_tested_positive_weekly"`
	Deaths                  string `json:"deaths_weekly"`
	Deaths_rate             string `json:"death_rate_weekly"`
}

type CCVIJsonRecords []struct {
	Geography_type      string `json:"geography_type"`
	Community_or_zip    string `json:"community_area_or_zip"`
	Community_area_name string `json:"community_area_name"`
	CCVI_score          string `json:"ccvi_score"`
	CCVI_category       string `json:"ccvi_category"`
}

type BuildingPermitsJsonRecords []struct {
	Id               string `json:"id"`
	Permit_nbr       string `json:"permit_"`
	Permit_type      string `json:"permit_type"`
	Review_type      string `json:"review_type"`
	Application_date string `json:"application_start_date"`
	Issue_date       string `json:"issue_date"`
	Total_fee        string `json:"total_fee"`
	Zip_code         string `json:"contact_1_zipcode"`
	Community_area   string `json:"community_area"`
}

type UnemploymentRatesJsonRecords []struct {
	Community_area      string `json:"community_area"`
	Community_area_name string `json:"community_area_name"`
	Unemployment        string `json:"unemployment"`
	Poverty             string `json:"below_poverty_level"`
	Income              string `json:"per_capita_income"`
}
