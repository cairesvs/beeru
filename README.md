# Beeru

Beeru is responsible for Get, Create and Find Points of sale aka "Pontos de venda" or shorthand PDV

![beeru](https://media.giphy.com/media/osbU9PXXgwHuM/giphy.gif)

> MMM

## Architecture

For the language I used golang with [gorilla mux](https://github.com/gorilla/mux) and for database [libpq port to golang](https://github.com/lib/pq).
[Logrus](github.com/sirupsen/logrus) to create structured log and docker for container.

## Documentation

Our packages have documentation to help us find useful functions and explanation about the functionalities.

### Run

```sh
make docs
```

To see the `beeru` documentation go to [Beeru Docs](http://localhost:6060/pkg/github.com/cairesvs/beeru/)

## How to Run

### Create Go dir:

```sh
mkdir $HOME/go
```

### Install Go Lang \o/:

If you are using mac, you can use brew:

```sh
brew install go
```

If you aren't, take a look at this page: https://golang.org/doc/install

Don't forget to add Go vars to your `.bash_profile`, `.zshrc` or whatever:

```sh
export GOPATH=$HOME/go
export GOROOT=/usr/local/opt/go/libexec  # This GOROOT is for MAC env
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
```

### Install Glide (Go package manager):

```sh
curl https://glide.sh/get | sh
```

### Postgres

We use postgres and postgis. If you are using mac you can use brew as well.

```sh
brew install postgres
brew install postgis
```

### Cloning this repo

We need to clone it in Go structure:

```sh
mkdir -p $HOME/go/src/github.com/cairesvs/
cd $HOME/go/src/github.com/cairesvs/
git clone git@github.com:cairesvs/beeru.git
cd beeru
```

### Install dependencies

```sh
make dependencies
```

And finally:

```sh
make run
```

## How to Test

Once all dependencies are satisfied you can run tests:

```sh
make test
```

## How to Deploy

We use docker to distribute our application. We have on `Makefile` three main tasks.

```sh
make docker/image DOCKER_IMAGE_VERSION=v0.0.1 
make docker/tag DOCKER_IMAGE_VERSION=v0.0.1 
make docker/push DOCKER_IMAGE_VERSION=v0.0.1 
```

TODO: We can use ECS, AWS with Cloudformation, k8s, gcloud or any docker supported provider.

## Under the hood

In this section core concepts of the application design will be brought to light, this is an effort to make
new contributors up to speed as fast as possible.

### Postgres persistence

For Get and Create PDVs we use simple postgres functionalities but for the `find` we use Postgis.

For more details about Postgis please read [this documentation](http://postgis.net/documentation/).

#### Find PDV give start point

```sql
SELECT id,trading_name, owner_name, document, ST_AsGeoJSON(coverage_area) as coverage_area, ST_AsGeoJSON(address) as address FROM pdv WHERE ST_Contains(pdv.coverage_area, ST_GeomFromText('POINT(%03.6f %03.6f)', 4326)) ORDER BY ST_Distance(pdv.coverage_area, pdv.address)
```

Since the coverage are is a MultiPolygon we use the `ST_Contains` to find if the given point intersects using the Postgis capabilities.