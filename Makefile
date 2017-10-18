test/create:
	curl -H "Content-Type: application/json" -X POST -d @test.json localhost:8000/pdv 2> /dev/null | jq

test/get:
	curl -X GET localhost:8000/pdv/16 2> /dev/null | jq

test/find:
	curl -X GET localhost:8000/pdvs/-43.297337/-23.013538 2> /dev/null | jq


postgres/up:
	docker run --name some-postgis -p 5433:5432 -e POSTGRES_PASSWORD=123 -d mdillon/postgis

postgres/down:
	docker rm -f some-postgis

postgres/connect:
	docker run -it --rm postgres psql -U postgres -h 127.0.0.1 -p 5433

run:
	go run main.go