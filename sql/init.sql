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