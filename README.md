# tracing-demo

Demo app for NAIS tracing based on Grafana stack.

This repository contains a `docker-compose` stack that sets up the observability stack locally on your computer.

Also doubles as a canary to test the Naiserator frontend integration.

## Local testing of logging stack

You can test your application using docker-compose:

```
docker-compose up
```

This will set up the following services:

- Grafana at http://localhost:3500, use it to explore your collected data
- Grafana Agent at `http://localhost:12347/collect`, use it to collect data from your frontend application
- Tempo, available on gRPC at `localhost:4317`, use it to collect traces from your backend
- Loki, only available inside Docker

## See also

- [NAIS documentation: Frontend application observability](https://docs.nais.io/observability/frontend/)
