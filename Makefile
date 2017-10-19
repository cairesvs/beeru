TESTS?=$$(go list ./... | egrep -v "vendor")
DOCKER_IMAGE_VERSION?=build
PG_NAME?=local-postgis
PG_PORT?=5433

test/create:
	curl -H "Content-Type: application/json" -X POST -d @test.json localhost:8000/pdv 2> /dev/null | jq

test/get:
	curl -X GET localhost:8000/pdv/16 2> /dev/null | jq

test/find:
	curl -X GET localhost:8000/pdvs/-43.297337/-23.013538 2> /dev/null | jq

postgres/up:
	docker run --name $(PG_NAME) -v ${PWD}/sql:/docker-entrypoint-initdb.d/ -p $(PG_PORT):5432 -d mdillon/postgis

postgres/down:
	docker rm -f $(PG_NAME)

postgres/connect:
	docker run -it --rm postgres psql -U postgres -h 127.0.0.1 -p $(PG_PORT)

dependencies:
	glide install

docs:
	godoc -http=":6060"

run:
	go run main.go


test:
	$(MAKE) postgres/up
	sleep 5
	PG_CONNECTION="user=postgres host=localhost dbname=postgres port=$(PG_PORT) sslmode=disable" go test -v $(TESTS)
	$(MAKE) postgres/down	

build:
	CGO_ENABLED=0 go build -v -a --installsuffix cgo --ldflags="-s" -o beeru

install:
	CGO_ENABLED=0 go install -v -a --installsuffix cgo --ldflags="-s"

docker/build:
	docker build -t caires/beeru:build -f Dockerfile.build .

docker/image: docker/build
	docker run --rm --entrypoint /bin/sh -v ${PWD}:/out:rw caires/beeru:build -c "cp /go/bin/beeru /out/beeru"
	docker build -t caires/beeru .

docker/tag:
	docker tag caires/beeru caires/beeru:${DOCKER_IMAGE_VERSION}

docker/run:
	docker run -ti -p 8000:8000 --rm caires/beeru:${DOCKER_IMAGE_VERSION}

docker/run/local:
	docker build -t caires/beeru:build -f Dockerfile.build .
	docker run --rm --entrypoint /bin/sh -v ${PWD}:/out:rw caires/beeru:build -c "cp /go/bin/beeru /out/beeru"
	docker build -t caires/beeru:local .
	docker run --name postgis -v ${PWD}/sql:/docker-entrypoint-initdb.d/ -e POSTGRES_PASSWORD=123 -d mdillon/postgis
	docker run -ti --link postgis:postgis -p 8000:8000 -e "PG_CONNECTION=user=postgres dbname=postgres sslmode=disable password=123 host=postgis" --rm caires/beeru:local

docker/stop/local:
	docker rm -f postgis

docker/push:
	docker push caires/beeru:${DOCKER_IMAGE_VERSION}
