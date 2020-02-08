# datadog-mock

[![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/nikscorp/datadog-mock)](https://hub.docker.com/repository/docker/nikscorp/rhymes)


## Overview

This is a fork from https://github.com/jancajthaml-devops/datadog-mock.
The purpose of fork is to modify the project structure and improve docker and compose files.

datadog-mock is a golang statsd mock server listening on port 8125 and relaying events to stdout.

### Changes
- simplify docker-compose file to only run the service
- add building the app and running autotests to Dockerfile
- DockerHub automated build
- switch to go dep (no experience about modules)
- add basic autotests with go testing

### Known issues
- incorrect handling of service checks and events dogstatd messeges (see skipped tests)

## Build and run tests

If you want to build image and run tests by yourself:
```
docker-compose build
```

## Run

```
docker-compose pull
docker-compose up -d
```

## License

This service is distributed under the Apache License, Version 2.0 license found
in the [LICENSE](./LICENSE) file.

Original project: https://github.com/jancajthaml-devops/datadog-mock