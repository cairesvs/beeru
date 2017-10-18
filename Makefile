TESTS?=$$(go list ./... | egrep -v "vendor")
DOCKER_IMAGE_VERSION?=build

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

dependencies:
	glide install

docs:
	godoc -http=":6060"

run:
	go run main.go

test:
	go test -v $(TESTS)

build:
	CGO_ENABLED=0 go build -v -a --installsuffix cgo --ldflags="-s" -o beeru

install:
	CGO_ENABLED=0 go install -v -a --installsuffix cgo --ldflags="-s"

docker/build:
	docker build -t cairesvs/beeru:build -f Dockerfile.build .

docker/image: docker_build
	docker run --rm --entrypoint /bin/sh -v ${PWD}:/out:rw cairesvs/beeru:build -c "cp /go/bin/beeru /out/beeru"
	docker build -t cairesvs/beeru .

docker/tag:
	docker tag cairesvs/beeru cairesvs/beeru:${DOCKER_IMAGE_VERSION}

docker/run:
	docker run -ti -p 8001:8000 --rm cairesvs/beeru:${DOCKER_IMAGE_VERSION}

docker/push:
	docker push cairesvs/beeru:${DOCKER_IMAGE_VERSION}
