WITH cct AS (
       SELECT
              DATE(
                     date_trunc('week', trip_start_timestamp) - interval '1 days'
              ) AS trip_week,
              pickup_community_area,
              dropoff_community_area,
              COUNT(DISTINCT trip_id) AS trip_cnt
       FROM
              taxi_trips AS tx
              INNER JOIN ccvi_ca AS ca ON tx.pickup_community_area = ca.community_area
              OR tx.dropoff_community_area = ca.community_area
       WHERE
              ca.ccvi_category = 'HIGH'
              AND tx.trip_start_timestamp >= '2020-1-1'
       GROUP BY
              trip_week,
              pickup_community_area,
              dropoff_community_area
)
SELECT
       fin.trip_week,
       can.community_area_name,
       SUM(fin.trip_cnt) AS trip_cnt
FROM
       (
              SELECT
                     trip_week,
                     pickup_community_area AS community_area,
                     SUM(trip_cnt) AS trip_cnt
              FROM
                     cct
              GROUP BY
                     trip_week,
                     community_area
              UNION
              SELECT
                     trip_week,
                     dropoff_community_area AS community_area,
                     SUM(trip_cnt) AS trip_cnt
              FROM
                     cct AS b
              GROUP BY
                     trip_week,
                     community_area
       ) AS fin
LEFT JOIN community_area_names AS can
ON fin.community_area = can.community_area
GROUP BY
       fin.trip_week,
       can.community_area_name
ORDER BY
       fin.trip_week,
       can.community_area_name;
