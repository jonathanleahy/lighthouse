# Backoffice Core BFF
![Badge](https://img.shields.io/badge/Go-v1.19-blue)
[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/pismo/backoffice-core-bff%2Fapi-deployment-integration?type=cf-1&key=eyJhbGciOiJIUzI1NiJ9.NjAxOTk1NjM1YzUyOWRiZDQzOTNlY2U5.ABCaPDSIUgw4UDE_tYjyLnqsuj85ISwyr1Pjd4jrWEo)]( https://g.codefresh.io/pipelines/edit/new/builds?id=6424866f79973f1bb0dca45d&pipeline=api-deployment-integration&projects=backoffice-core-bff&projectId=6424866d7afcb316d29820eb)
<a href="https://sonar.pismo.services/dashboard?id=backoffice-core-bff" target="_blank">![Badge](https://img.shields.io/badge/Sonar-Code%20Quality-informational?&logo=sonarqube)</a>

This project aims to be a backend for frontend for the backoffice service.

## Architecture
Project was based on "clean architecture". For more details, access [here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

## Project Configuration
To deploy your application in production environments, you **must** configure these following environment variables:

| Environment Variable           | Default Value                                    | Description |
|--------------------------------|--------------------------------------------------|-------------|
| APP_NAME                       | backoffice-core-bff                              |             |
| VERSION                        | unset                                            |             |
| ENV                            | dev                                              |             |
| LOG_LEVEL                      | info                                             |             |
| OTEL_EXPORTER_OTLP_ENDPOINT_GO | localhost:4317                                   |             |
| PORT                           | 8080                                             |             |
| SERVER_VERBOSE                 | false                                            |             |
| HTTP_DEFAULT_TIMEOUT           | 25                                               | In Seconds  |
| AWS_REGION                     | sa-east-1                                        |             |
| SNS_CONSOLE_AUDIT              | arn:aws:sns:sa-east-1:270036487593:console-audit |             |
| CORS_ENABLED                   | true                                             |             |
| CORS_ALLOWED_ORIGINS           | *                                                |             |
| CORS_ALLOWED_METHODS           | GET,HEAD,PUT,PATCH,POST,DELETE                   |             |
| CORS_ALLOWED_HEADERS           | *                                                |             |
| CORS_ALLOWED_CR3D3NT1ALS       | false                                            |             |
| CORS_EXPOSED_HEADERS           |                                                  |             |
| CORS_MAX_AGE                   | 86400                                            | In seconds  |


## Running the Application

### Run BFF
make run-api
### Run Tests
make test
### Update OpenApi
To generate openapi files, we need to run this following command
make openapi-gen
By default, documents are generated in the docs/openapi directory.
### Update Mocks
To generate Mocks files, we need to run this following command
make build-mock
### Makefile
Use ``make help`` or only ``make`` to check all of the available commands.

## Endpoints

| Endpoint  | Method | Auth? | Description                  | Documentation                             |
|-----------|--------|-------|------------------------------|-------------------------------------------|
| `/health` | GET    | Yes   | Check the API current status |                                           |
| `/query`  | POST   | Yes   | Graphql queries              |                                           |


## Documentation
1. [Swagger](docs/openapi/swagger.json)
2. [Dashboard](https://pismo.grafana.net/d/xBr3AK97z/main-dashboard?orgId=1&refresh=5s&var-container=backoffice-core-bff&var-cid=.%2A&var-dataset=sandbox&var-environment=integration)

## About us
This project is maintained by **squad-enablement**, please feel free to contact.
- Slack channel: [#squad-enablement](https://pismo.slack.com/archives/C04M33MSC5P)
