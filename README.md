[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/b3ntly/twelvefactor_databases/master/LICENSE.txt) 
[![Build Status](https://travis-ci.org/b3ntly/twelvefactor_databases.svg?branch=master)](https://travis-ci.org/b3ntly/twelvefactor_databases)
[![Go Report Card](https://goreportcard.com/badge/github.com/b3ntly/twelvefactor_databases)](https://goreportcard.com/report/github.com/b3ntly/twelvefactor_databases)

### Sample 12-Factor Application with Golang

This application demonstrates database usage within a Golang application using Postgres.

### Usage

```bash
docker run gcr.io/twelvefactor/twelvefactor_databases
```

### Environment Options

These are configuration variables that can be passed to the docker container.

| Key | Description | Default |
| ------------- |:-------------:| -----:|
| PORT | The port on 127.0.0.1 from which this application will serve. | 9090 |
| PING_PATH | The URL path at which to serve responses | /ping |
| PING_RESPONSE | The string response returned by a GET request to /ping | PONG |
| REQ_TIMEOUT | Request Timeout in Milliseconds | 500 |
| SERVER_READ_TIMEOUT | Server Read Timeout in Milliseconds | 1000 |
| SERVER_WRITE_TIMEOUT | Server Write Timeout in Milliseconds | 2000 |
| DB_CONN_MAX_LIFETIME | Max duration of a database timeout | unlimited |
| DB_MAX_OPEN | Max open database connections | unlimited |
| DB_MAX_IDLE | Max idle database connections | 2 |
| POSTGRES_URI | Database connection URI | postgresql://postgres@localhost:5432/postgres?sslmode=disable |
| USERS_PATH | Path to expose the users service | /users |
| USERS_SELECT_LIMIT | The number of users to return from a GET request to the USERS_PATH | 10 |