# Computing Cultural Heritage in the Cloud

## America's Public Bible: Machine-Learning Detection of Biblical Quotations across Library of Congress Collections via Cloud Computing

[Lincoln Mullen](https://lincolnmullen.com), Director of Computational History,
[Roy Rosenzweig Center for History and New Media](https://rrchnm.org), George
Mason University

### Project status

[![Go](https://github.com/lmullen/cchc/actions/workflows/go.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/go.yml)
[![Docker](https://github.com/lmullen/cchc/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/lmullen/cchc/actions/workflows/docker-publish.yml)

### About this repository

This repository contains code for one of the projects that are part of the [Computing Cultural Heritage in the Cloud](https://labs.loc.gov/work/experiments/cchc/) initiative at the Library of Congress Labs.

### Using this repository

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

You can run the application with `make run`.

Note that the application is containerized, except for the database.

### License

All code is copyrighted &copy; 2021 Lincoln A. Mullen. Code is licensed [CC0 1.0
Universal](https://github.com/lmullen/cchc/blob/main/LICENSE).

