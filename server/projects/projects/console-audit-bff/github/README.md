<h1 align="center">🚀 Console Audit BFF 🚀</h1><div align="center">

![Badge](https://img.shields.io/badge/Go-v1.19-blue)
![Badge](https://img.shields.io/badge/Status-Production%20Ready-green)
<a href="https://pismo.grafana.net/d/xBr3AK97z/main-dashboard?orgId=1&refresh=5s&var-container=console-audit-bff&var-cid=.%2A&var-dataset=sandbox&var-environment=integration" target="_blank">![Badge](https://img.shields.io/badge/Grafana%20Dashboard-yellow?&logo=grafana)</a>
<a href="https://app.getcortexapp.com/admin/service/25341/homepage" target="_blank">![Badge](https://img.shields.io/badge/Cortex-console--audit--bff-green)</a>
<a href="https://sonar.pismo.services/dashboard?branch=main&id=console-audit-bff" target="_blank">![Badge](https://img.shields.io/badge/Sonar-Code%20Quality-informational?&logo=sonarqube)</a>
<a href="https://pismo.slack.com/archives/C031SBDMEEA" target="_blank">![Badge](https://img.shields.io/badge/slack-%23squad--console--core-informational?style=social&logo=slack)</a>

</div>

---

<p align="center">
 • <a href="#-about-the-console-audit-bff">About the Console Audit BFF</a> •
 <a href="#-project-configuration">Project Configuration</a> •
 <a href="#-running-the-application">Running the Application</a> •
</p>
<p align="center">
 • <a href="#-endpoints">Endpoints</a> •
 <a href="#-production-ready-integrations">Production Ready Integrations</a> •
</p>
<p align="center">
 • <a href="#-documentation">Documentation</a> •
 <a href="#-about-us">About Us</a> •
</p>

---

## 💻 About the Console Audit BFF

The console-audit-bff is responsible for servers specifics endpoints console audit domain.

This project is maintained by <a href="https://pismo.slack.com/archives/C031SBDMEEA" target="_blank">
#squad-console-core</a>.

---

## 🛠 Project Configuration

If you don't have psm-sdk installed, install it
following [How to Install the Pismo SDK](https://github.com/pismo/psm-sdk/#-installation).

After run the command below:

go mod tidy

⚠️ To deploy your application in production environments, you **must** configure these following environment variables:

| Environment Variable           | Default Value                                          | Description |
|--------------------------------|--------------------------------------------------------|-------------|
| APP_NAME                       | console-audit-bff                                      |             |
| VERSION                        | unset                                                  |             |
| ENV                            | dev                                                    |             |
| LOG_LEVEL                      | info                                                   |             |
| OTEL_EXPORTER_OTLP_ENDPOINT_GO | localhost:4317                                         |             |
| PORT                           | 8080                                                   |             |
| SERVER_VERBOSE                 | false                                                  |             |
| HTTP_DEFAULT_TIMEOUT           | 25                                                     | In Seconds  |
| AWS_REGION                     | sa-east-1                                              |             |
| SNS_CONSOLE_AUDIT              | arn:aws:sns:sa-east-1:270036487593:console-audit       |             |
| CORS_ENABLED                   | true                                                   |             |
| CORS_ALLOWED_ORIGINS           | *                                                      |             |
| CORS_ALLOWED_METHODS           | GET,HEAD,PUT,PATCH,POST,DELETE                         |             |
| CORS_ALLOWED_HEADERS           | *                                                      |             |
| CORS_ALLOWED_CR3D3NT1ALS       | false                                                  |             |
| CORS_EXPOSED_HEADERS           |                                                        |             |
| CORS_MAX_AGE                   | 86400                                                  | In seconds  |
 | CONSOLE_AUDIT_API_URL          | https://console-audit-api.integration.pismolabs.io:443 |             |

Configuration example for IntelliJ IDEA:

### Local example

`AWS_REGION=sa-east-1;AWS_PROFILE=pismolabs-ext;LOG_LEVEL=debug`

### EXT / Integration example

`AWS_REGION=sa-east-1;AWS_PROFILE=pismolabs-ext;LOG_LEVEL=debug`

## 🚀 Running the Application

### Run BFF

make docker-compose-up
make run-api

### Run Tests

make docker-compose-up
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

## 📚 Endpoints

| Endpoint | Method | Auth? | Description      | Documentation                             |
|----------|--------|-------|------------------|-------------------------------------------|
| `/query` | POST   | Yes   | GraphQL Endpoint | [Link](docs/query/query.findAuditByID.MD) |

## 🎯 Production Ready Integrations

All our features are 100% integrated with [Console Audit](https://console.pismo.io/), Grafana, Honeycomb and AWS.

Audit list in console:
![Audit list in console](docs/images/console-audit-list-example.png)

Audit detail in console:
![Audit detail in console](docs/images/console-audit-detail-example.png)

## 📚 Documentation

1. [How to install the Pismo SDK](https://github.com/pismo/psm-sdk/#-installation)
2. [Docs](docs)
3. [Swagger](docs/openapi/swagger.json)
4. [Dashboard](https://pismo.grafana.net/d/xBr3AK97z/main-dashboard?orgId=1&refresh=5s&var-container=console-audit-bff&var-cid=.%2A&var-dataset=sandbox&var-environment=integration)

## 📚 About us

This project is maintained by **squad-console-core**, please feel free to contact.

- Slack channel: [#squad-console-core](https://pismo.slack.com/archives/C031SBDMEEA)

