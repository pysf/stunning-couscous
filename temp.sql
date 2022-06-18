create extension if not exists cube; create extension if not exists earthdistance;


select earth_distance(
  ll_to_earth(52.528971849007036 ,13.430548464498173),
  ll_to_earth(52.51999140,13.40497255)
); 

CREATE INDEX CONCURRENTLY partnerloc
    ON partner USING gist (ll_to_earth("location"[0],"location"[1]));

DROP INDEX partnerloc;

SELECT COUNT(*) FROM partner;

SELECT id, distance,location,rating ,operatingradius  FROM (

select id,location,operatingradius, rating,earth_distance(
  ll_to_earth("location"[0],"location"[1]),
  ll_to_earth(52.51999140,13.40497255)
) as distance FROM partner

) a WHERE operatingradius > distance ORDER by rating DESC,distance ASC ; 


-- With INDEX
-- Gather Merge  (cost=65840.41..72499.28 rows=57072 width=36)
--   Workers Planned: 2
--   ->  Sort  (cost=64840.39..64911.73 rows=28536 width=36)
-- "        Sort Key: partner.rating DESC, (sec_to_gc(cube_distance((ll_to_earth(partner.location[0], partner.location[1]))::cube, '(3775281.5188042056, 899745.1724846419, 5061495.343728965)'::cube)))"
--         ->  Parallel Seq Scan on partner  (cost=0.00..62728.65 rows=28536 width=36)
-- "              Filter: ((operatingradius)::double precision > sec_to_gc(cube_distance((ll_to_earth(location[0], location[1]))::cube, '(3775281.5188042056, 899745.1724846419, 5061495.343728965)'::cube)))"

