[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/b3ntly/twelvefactor_ping/master/LICENSE.txt) 
[![Build Status](https://travis-ci.org/b3ntly/twelvefactor_ping.svg?branch=master)](https://travis-ci.org/b3ntly/twelvefactor_ping)
[![Go Report Card](https://goreportcard.com/badge/github.com/b3ntly/twelvefactor_ping)](https://goreportcard.com/report/github.com/b3ntly/twelvefactor_ping)

### Sample 12-Factor Application with Golang

This application returns simple ping/pong style responses as a demonstration
of the 12-factor application methodology.

### Usage

```bash
docker run gcr.io/twelvefactor/twelvefactor_ping
```

### Environment Options

These are configuration variables that can be passed to the docker container.

| Key | Description | Default |
| ------------- |:-------------:| -----:|
| PORT | The port on 127.0.0.1 from which this application will serve. | 9090 |
| ENDPOINT | The URL path at which to serve responses | /ping |
| RESPONSE | The string response returned by a GET request to /ping | PONG |
| REQ_TIMEOUT | Request Timeout in Milliseconds | 500 |
| SERVER_READ_TIMEOUT | Server Read Timeout in Milliseconds | 1000 |
| SERVER_WRITE_TIMEOUT | Server Write Timeout in Milliseconds | 2000 |