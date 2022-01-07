# Computing Cultural Heritage in the Cloud

## America's Public Bible: Machine-Learning Detection of Biblical Quotations across Library of Congress Collections via Cloud Computing

[Lincoln Mullen](https://lincolnmullen.com), Director of Computational History,
[Roy Rosenzweig Center for History and New Media](https://rrchnm.org), George
Mason University

[![cchc-crawler container](https://github.com/lmullen/cchc/actions/workflows/docker-publish-crawler.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/docker-publish-crawler.yml)

[![cchc-itemmd container](https://github.com/lmullen/cchc/actions/workflows/docker-publish-itemmd.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/docker-publish-itemmd.yml)

[![cchc-language-detector container](https://github.com/lmullen/cchc/actions/workflows/docker-publish-language-detector.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/docker-publish-language-detector.yml)

[![cchc-predictor container](https://github.com/lmullen/cchc/actions/workflows/docker-publish-predictor.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/docker-publish-predictor.yml)

[![cchc-ctrl container](https://github.com/lmullen/cchc/actions/workflows/docker-publish-ctrl.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/docker-publish-ctrl.yml)

[![Go tests](https://github.com/lmullen/cchc/actions/workflows/go.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/go.yml) 

### About this repository

This repository contains code for one of the projects that are part of the [Computing Cultural Heritage in the Cloud](https://labs.loc.gov/work/experiments/cchc/) initiative at the Library of Congress Labs.

The purpose of this project is to experiment with cloud computing approaches to computational history, using the Library of Congress's [digital collections](https://www.loc.gov/collections/) (specifically, the full text collections). This experiment has three parts. First, it uses the [loc.gov APIs](https://www.loc.gov/apis/) whenever possible. Rather than use a batch processing model, it can be run continuously to receive updates (and process them) as the Library of Congress adds new materials. Second, this application is containerized and split apart into several microservices so that it can be run on a variety of hardware (e.g., a laptop, a server, cloud compute, a high-performance cluster). And third, it allows the user to both crawl the Library of Congress digital collections and write jobs that can do work on the collection.  There are two included job processors: one to look for multilingual documents in the Library of Congress and the other to look for biblical quotations (as in [*America's Public Bible](https://americaspublicbible.org)).

### License

All code is copyrighted &copy; 2021 Lincoln A. Mullen. Code is licensed [CC0 1.0
Universal](https://github.com/lmullen/cchc/blob/main/LICENSE).

## Using the application

This application has a few main sections:

- A PostgreSQL database, into which the metadata from the Library of Congress API is stored, along with the results of running jobs on the collections.
- A crawler and item metadata fetcher which get information from the loc.gov API.
- Services which run jobs (i.e., do useful work) on the collection.
- Some utilities for managing the application state.

All of these parts of the application are containerized, though you are strongly encouraged to use your own, non-containerized database.

When using Docker Compose for orchestration, the application is split up into several profiles which do discrete kinds of work. These are the profiles (specific to Docker Compose), as well as the description of each of the services/containers associated with them.

- Compose profile: `db`. An optional profile to set up a database
	- Compose service: `db`; Docker image: `postgres:14`. A PostgreSQL database
	- Compose service: `adminer`; Docker image: `adminer`. A database management client
- Compose profile: `api`. Fetches metadata from the loc.gov API
	- Compose service: `crawler`; Docker image: `ghcr.io/lmullen/cchc-crawler`. Discovers items from the API.
	- Compose service: `itemmd`; Docker image: `ghcr.io/lmullen/cchc-itemmd`. Fetches the full metadata for items that have been discovered by the crawler.
- Compose profile: `ctrl`
	- Compose service: `ctrl`; Docker image: `ghcr.io/lmullen/cchc-ctrl`. Runs administrative scripts on the database.
- Compose profile: `languages`
	- Compose service: `language-detector`; Docker image: `ghcr.io/lmullen/cchc-language-detector`. Guess the language of each sentence in full-text items, to look for multilingual documents.
- Compose profile:  `quotations`
	- Compose service: `predictor`; Docker image: `ghcr.io/lmullen/cchc-predictor`. Identifies biblical quotations in full-text items.

If you are not using Docker Compose to run these containers, you can still use this organization conceptually to understand the pieces of the application. See below for a description of each of the services.

### Settings

Application-wide settings are set with environment variables.

- `CCHC_DBSTR`: This is the URL to the PostgreSQL database. This setting is not optional, and each service will fail without it. It should take the following form: `postgres://user:password@hostname:5432/database?sslmode=disable`
- `CCHC_LOGLEVEL`: This is an optional setting to control the verbosity of logging. You can set it to any of the following values: `error`, `warn`, `info`, `debug`. The default level is `info`.
- `CCHC_VERSION`: This is an optional setting to control which version of the containers you get when pulling them. By default, the value used will be `release`, which will get only known-good releases. Only set or change this if you really know that you should. Setting this to `main` will get containers from the `main` branch.

### Docker

You can run each of the parts of this application using [Docker](https://www.docker.com). The Docker images for the different services are provided in the [GitHub container registry](https://github.com/lmullen/cchc/packages). For instance, you could run all of the services on a single machine using `docker compose`, as explained below. But you could also run the crawler and item metadata fetcher continuously in the cloud, then scale up your job processors on a high-performance cluster periodically. (Details of using Docker for deploying this software are beyond the scope of this project.)

[Docker Compose](https://docs.docker.com/compose/) is a way of orchestrating many containers. This repository contains a [`docker-compose.yml`](https://github.com/lmullen/cchc/blob/main/docker-compose.yml) file which defines the relationship between the containers. The remainder of this documentation assumes you are using Docker Compose (v2.2 or greater) to run the application. This comes standard with recent versions of [Docker Desktop](https://www.docker.com/products/docker-desktop), and it can also be installed on servers. You should be able to run the command `docker compose version` to check your installation.

### Cloning the repository / pulling the containers

You can get this application's source code by [cloning the repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository). 

When you navigate in your terminal to the root of the repository, you can get run the necessary commands using Docker Compose. Use the following command to pull all the necessary containers:

```
docker compose --profile db --profile cchc pull
```

### PostgreSQL database

All metadata, application state (such as the job queues) and results are stored in a [PostgreSQL database](https://www.postgresql.org). While any version from v12 up should work, this application has been tested the most on PostgreSQL 14.1.

You are **strongly** encouraged to use a database whose data is stored in a way that you can understand and access if you care at all about the data stored in it. While this application does provide a Docker container which stores its data in a [Docker volume](https://docs.docker.com/storage/volumes/), persisting data in Docker requires considerable understanding to do correctly, as does securing the database.

Your database should have a user with permissions to create schemas and tables. The actual migrations will be handled by `cchc-ctrl` service (see below). You can pass the database connection information into the application using the `CCHC_DBSTR` environment variable as detailed [above](#settings).

If you would like to use a containerized database, then you can start one with Docker Compose.

```
docker compose --profile db up --detach
```

If you set the `CCHC_DB_PASSWORD` environment variable _before_ the first time you run the database service, it will set the PostgreSQL user password to that value. Otherwise, it will use the default password below.

These are the configuration settings for the database:

- User name: `cchc`
- Database name: `cchc`
- Password: `cchcdbpass`
- Host: `localhost` (assuming you are running this on the same machine)
- Port: `5433` (a non-default port to avoid conflicting with other PostgreSQL instances)

If you are using this containerized database, then you can set the database connection string to an environment variable as follows:

```
export CCHC_DBSTR=postgres://cchc:cchcdbpass@localhost:5433/cchc?sslmode=disable
```

Otherwise, use that connection string as a template.

You can use whatever client you prefer to connect to the database. However, if you use the containerized database, then a [basic client](https://www.adminer.org) is started for you. You can access this service at <http://localhost:8888/>. Use the following settings to login:

![Adminer login](README-adminer-login.png)

### Initializing the database and using `cchc-ctrl`

Regardless of whether you use the containerized database or host your own, you will need to migrate the database to the current schema. You can do this, as well as other maintenance tasks, using the `ctrl` service.

First check that you can successfully connect to the database.

```
docker compose run --rm ctrl /cchc-ctrl ping
```

Then you can run the migrations:

```
docker compose run --rm ctrl /cchc-ctrl migrate
```

Currently, this utility supports the following actions:

- `help`:        Help about any command
- `migrate`:     Migrate the database to the current schema
- `ping`:        Check connection to the database
- `reset`:       Reset the database (deletes all data)
- `retry-jobs`   Retry skipped and failed jobs

For full documentation on how to use this utility, consult the help.

```
docker compose run --rm ctrl /cchc-ctrl help
```

### Crawler and item metadata fetcher

Two services (`crawler` and `itemmd`) identify items from the Library of Congress API and then fetch the full metadata. These services are intended to be run continuously. The crawler will periodically (currently, once every two days) check for updates to the Library of Congress digital collections, and the item metadata fetcher will get the full metadata for each item.

They save the resulting metadata in several database tables in the `public` schema, including `collections` (digital collections from LOC), `items` (specific items, which are associated with one or more collections), and `resources` and `files`, which track the files associated with items. The `api` column on the `items` table contains the full JSON response for each item from the API, and can be used to get other metadata fields which have not been extracted into specific columns.

To start these services, run the following:

```
docker compose --profile api up --detach
```

Note that in the documentation below the `--scale` flag is suggested for using more than one worker at a time. Do not attempt to scale the crawler or item metadata fetcher beyond one instance each. The loc.gov API is strictly rate limited, and using more than one replica for these services will result in your IP address being blocked.

### Language detector

This service seeks to identify the language of each sentence in the full-text items, and thus identify multilingual documents in the collections.

To start this service, run the following:

```
docker compose --profile langauges up --detach
```

Results as stored in the `results.languages` table. This service keeps track of jobs in the `jobs.fulltext` table. This computed result can then be compared to the `language` field in the `items` table. Items that do not have full text will be skipped. You can delete skipped or failed jobs with the `cchc-ctrl` service. 

This is an example of a service which does useful work on the Library of Congress collections. Note that this particular service is written entirely in Go. It could be used as a template for creating a similar service in Go.

### Miscellaneous

To stop and remove a particular service, you can use the `stop` or `down` functions in Docker compose. To stop and remove all services (including the database), run the following:

```
docker compose --profile db --profile cchc down
```

While containers are running or stopped (but not after they have been removed) you can see the logs by running the following:

```
docker compose logs
```

You can see the logs for a particular service by running the following, changing the name for the service you want.

```
docker compose logs itemmd
```
