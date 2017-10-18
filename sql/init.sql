CREATE EXTENSION postgis;

CREATE TABLE pdv
(
  id SERIAL NOT NULL PRIMARY KEY,
  trading_name VARCHAR(255),
  owner_name VARCHAR(255),
  document CHAR(17) UNIQUE,
  coverage_area geometry(MultiPolygon,4326),
  address geometry(Point,4326)
);

SELECT pdv.id FROM pdv WHERE ST_Contains(pdv.coverage_area, ST_GeomFromText('POINT(-43.397337 -23.213538)', 4326)) ORDER BY ST_Distance(pdv.coverage_area, pdv.address); 