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

### Settings

Application-wide settings are set with environment variables.

- `CCHC_DBSTR`: This is the URL to the PostgreSQL database. This setting is not optional, and each service will fail without it. It should take the following form: `postgres://user:password@hostname:5432/database?sslmode=disable`
- `CCHC_LOGLEVEL`: This is an optional setting to control the verbosity of logging. You can set it to any of the following values: `error`, `warn`, `info`, `debug`. The default level is `info`.

### Dependencies

#### Docker

You can run each of the parts of this application using [Docker](https://www.docker.com). The Docker images for the different services are provided in the [GitHub container registry](https://github.com/lmullen/cchc/packages). For instance, you could run 

#### PostgreSQL database

This application assumes that configuration is passed in as environment
variables. You should set the following environment variables, though most will
have reasonable defaults.

Database (PostgreSQL) configuration:

- `CCHC_DBHOST`
- `CCHC_DBPORT`
- `CCHC_DBUSER`
- `CCHC_DBPASS`
- `CCHC_DBNAME`

Message broker (RabbitMQ) configuration:

- `CCHC_QUSER`
- `CCHC_QPORT`
- `CCHC_QHOST`
- `CCHC_QPASS`

Application configuration:

- `CCHC_LOGLEVEL`

The `Makefile` controls most of the application. You can create the database
with `make db-create` and run migrations with `make db-up` and `make db-down`.

You can run the application with `make up`.

Note that the application is containerized, except for the database. You will
need to provide your own PostgreSQL database.


