INSERT INTO community_area_names (community_area, community_area_name)
SELECT  DISTINCT community_area
       ,community_area_name
FROM unemployment
ORDER BY community_area;
