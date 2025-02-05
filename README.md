# Galadriel

[![CodeQL](https://github.com/HewlettPackard/galadriel/actions/workflows/codeql.yml/badge.svg)](https://github.com/HewlettPackard/galadriel/actions/workflows/codeql.yml)
[![PR Build](https://github.com/HewlettPackard/galadriel/actions/workflows/linter.yml/badge.svg)](https://github.com/HewlettPackard/galadriel/actions/workflows/linter.yml)
[![Scorecards supply-chain security](https://github.com/HewlettPackard/galadriel/actions/workflows/scorecards.yml/badge.svg)](https://github.com/HewlettPackard/galadriel/actions/workflows/scorecards.yml)
[![trivy](https://github.com/HewlettPackard/galadriel/actions/workflows/trivy.yml/badge.svg)](https://github.com/HewlettPackard/galadriel/actions/workflows/trivy.yml)

## Development

### REST API

Server and Client Go code is generated from the OpenAPI definition by [oapi-codegen](https://github.com/deepmap/oapi-codegen).

Run the following command to generate the code, from the root folder:

```bash
oapi-codegen --config=spec/api/schemas.cfg.yaml spec/api/schemas.yaml
oapi-codegen --config=spec/api/harvester.cfg.yaml spec/api/harvester.yaml
oapi-codegen --config=spec/api/management.cfg.yaml spec/api/management.yaml
```

Run the following command to have a live view of the API documentation:

```bash
make api-doc
```
This will grab the OpenAPI `api.yaml` file and generate a website for exploring and testing the API. Further changes in your API definition file can be reloaded by refreshing the website (it may require you to hard-refresh to avoid caching issues). The REST API documentation will be available in your browser at `http://localhost:8000`.

### Testing

There are a few make targets available to test the code:

* `make test`: Runs all tests.
* `make coverage`: Runs all unit tests and reports back test coverage. More details can be found in the file `./out/coverage/index.html`.
