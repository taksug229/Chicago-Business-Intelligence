SELECT
       DATE(
              date_trunc('week', trip_start_timestamp) - interval '1 days'
       ) AS trip_week,
       dropoff_zip_code,
       COUNT(DISTINCT trip_id) AS trip_cnt
FROM
       taxi_trips
WHERE
       pickup_community_area IN (76, 56, 64)
       AND trip_start_timestamp >= '2020-1-1'
GROUP BY
       trip_week,
       dropoff_zip_code
ORDER BY
       trip_week,
       dropoff_zip_code;
