WITH ca AS (
       SELECT
              community_area
       FROM
              unemployment
       WHERE
              per_capita_income < 30000
),
zip AS (
       SELECT
              DISTINCT zci.zip_code
       FROM
              (
                     SELECT
                            DISTINCT zc.zip_code,
                            zc.community_area
                     FROM
                            (
                                   SELECT
                                          zip_code,
                                          community_area1 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area2 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area3 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area4 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area5 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area6 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area7 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area8 AS community_area
                                   FROM
                                          zip2ca
                                   UNION
                                   SELECT
                                          zip_code,
                                          community_area9 AS community_area
                                   FROM
                                          zip2ca
                            ) AS zc
                     WHERE
                            zc.community_area IS NOT NULL
              ) AS zci
              INNER JOIN ca ON zci.community_area = ca.community_area
)
SELECT
       EXTRACT(
              YEAR
              FROM
                     bp.issue_date
       ) AS issue_year,
       bp.zip_code,
       SUM(bp.total_fee) AS total_fee
FROM
       (
              SELECT
                     issue_date,
                     zip_code,
                     total_fee
              FROM
                     building_permits
              WHERE
                     permit_type = 'PERMIT - NEW CONSTRUCTION'
                     AND issue_date >= '2020-1-1'
       ) AS bp
       INNER JOIN zip ON bp.zip_code = zip.zip_code
GROUP BY
       issue_year,
       bp.zip_code
ORDER BY
       issue_year,
       bp.zip_code
