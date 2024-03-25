WITH rate AS (
       SELECT
              DISTINCT upr.community_area,
              upr.community_area_name,
              upr.unemployment_rate,
              upr.poverty_rate
       FROM
              (
                     (
                            SELECT
                                   community_area,
                                   community_area_name,
                                   unemployment_rate,
                                   poverty_rate
                            FROM
                                   unemployment
                            ORDER BY
                                   unemployment_rate DESC
                            LIMIT
                                   5
                     )
                     UNION
                     (
                            SELECT
                                   community_area,
                                   community_area_name,
                                   unemployment_rate,
                                   poverty_rate
                            FROM
                                   unemployment
                            ORDER BY
                                   poverty_rate DESC
                            LIMIT
                                   5
                     )
              ) AS upr
)
SELECT
       EXTRACT(
              YEAR
              FROM
                     bp.issue_date
       ) AS issue_year,
       r.community_area_name,
       SUM(bp.total_fee) AS total_fee
FROM
       building_permits AS bp
       INNER JOIN rate AS r ON bp.community_area = r.community_area
WHERE bp.issue_date >= '2020-1-1'
GROUP BY
       issue_year,
       r.community_area_name
