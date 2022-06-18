create extension if not exists cube;
create extension if not exists earthdistance;


select earth_distance(
  ll_to_earth(52.528971849007036 ,13.430548464498173),
  ll_to_earth(52.51999140,13.40497255)
); 


SELECT a.id, a.distance,a.location ,a.operatingradius  FROM (

select id,location,operatingradius, earth_distance(
  ll_to_earth("location"[0],"location"[1]),
  ll_to_earth(47.237866887419074,13.404972550000002)
) as distance FROM partner

) a WHERE operatingradius > distance; 
