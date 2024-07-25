# FX Rate Getter API Service

This REST service provides highly available foreign exchange rates for onr or more quote currencies, 
by using multiple provider API sources.

Results may be cached for a set duration, which reduces third-party service load. 
It uses configurable load balancing and aggregation strategies to ensure consistent, fast 
and potentially free* FX rate responses, minimizing downtime risk and distributing load across providers.

This app can be adapted from being domain-specific (an FX rate retriever) to being a general data retriever 
for any data type from multiple similar APIs, such as weather or stock prices. Contributions welcome.

## Features:
- Configurable (see `config.json`)
- Multiple FX rate API provider sources (extensible to add more)
- 6 [load balancing and routing strategies](#load-balancing-and-routing-strategies) to fetch rates from healthy providers
- Rate limiter for requests, uses a `fixed bucket` algorithm
- Error handling with unique codes, formatted tracing, console, and logging output
- Caching for rates with configurable expiry
- Allow-list for supported currencies for your service
- Collects operational statistics
- Healthcheck endpoint to monitor the service and its providers

## API Endpoints:
```http
GET /rates?from={aaa}&to={bbb,ccc}    Eg: /rates?from=USD&to=EUR,USD,JPY
GET /rate/{from}/{to}                 Eg: /rate/USD/EUR
GET /status
GET /health
```

## How to run:
1. Clone the repository and download dependencies.
2. Set up your [config file](#setting-up-the-config-file) (`config.json`).
3. Run the application:
    - Go run:
        - Run `go run main.go` to start the server.
    - Makefile:
        - Run `make help` to see available commands.
        - Run `make build-run` to build and start the server.
    - Docker Compose:
        - Run `docker-compose up` to start the server in a container. 
4. To run tests:
    - Run `make test`, or `go test ./...` to run unit tests.
    - Run `make test-cover`, or `go test -cover ./...` to run tests with coverage.

Run `$ make help` for a full list of available commands.

### Setting up the config file:
- Rename `config.example.json` to **config.json**.
- Add your **API keys** for the providers you want to use and **enable** them.
- Set the **load balancing strategy** you want to use.
- Set the **cache duration**.
- Set the **rate limiter** configuration.
- Set your enabled **currencies**.
- Select the **router** you want to use (`gin` or `fiber`).
- Set the **port** you want to run the server on.

### Supported adapters*:
- Currency Layer (apilayer.com)
- Fixer API (apilayer.com)
- ExchangeRate-API (exchangerate-api.com)
- Free Currency API (freecurrencyapi.com)
- Free Currency Converter (currencyconverterapi.com)
- Open Exchange Rates (openexchangerates.org)

> **Note:** Ensure you read and comply with the terms and conditions of the third-party API providers you choose to use.

### Load balancing and routing strategies:
- **First** (Default):
    - Fetches rates from the first provider in the list.
    - If the provider is not healthy, it will try the next one.
- **Aggregate**:
    - Average of all the rates from all the enabled and healthy providers.
    - Disabled or unhealthy providers are ignored.
- **Priority**:
    - Fetches rates from healthy providers in a priority order, specified in the config.
- **Race**:
    - Concurrently fetches rates from all providers and returns the fastest response.
- **Random**:
    - Fetches rates from a random healthy provider.
- **Round Robin**:
    - Iterates through the list of healthy providers in a round-robin fashion.
    - Selects a different provider for each request.

The internal cache, when enabled, is always preferred regardless of the load balancing strategy. Where results have been aggregated from different providers, the cache will store the mean rate.

> **Note:** The application is easily extensible to support more providers and load balancing strategies.

### Application architecture:
- Router agnostic design, supports both `Gin` and `Fiber` as configurable. Easily add your preferred router.
- In-memory cache to store the most recent rates. The cache is extensible to other drivers (Redis, AWS, etc.).
- Error bundle to handle errors with unique codes, messages, printing, and chaining.
- Nicely formatted console output with colors.
- Panic recovery to catch panics and continue running.
- Pass a CLI argument `--config` to specify a custom config file, or it will look for `config.json` in the root directory.
- Environment arguments (for supported keys/values) will override the config file arguments, as configurable.

### Issues
If you encounter any issues, please open an issue on the repository.

#### Known issues:
- Some files are not fully covered by tests.

### Want to contribute? Possible improvements include:
- Add a back-off strategy per provider so that if they get rate limited, they will pause for a while before trying again.
- Add basic-auth or token-based authorization for administrative endpoints like `/status`, etc.
- Endpoint to refresh provider initializations:
    - For the ones that did not previously start successfully or
    - After a while in case their list of supported currencies changed.
- Stats should collect the number of times each provider was hit.
- Stats should compute and save the number of API calls per minute, hour, day.
- Option to support JSON RPC for the API.
- Option to support gRPC for the API.
- Option to support GraphQL for the API.
- Support a database of historic rates:
    - Expose endpoints to query them.
    - Provide options to select database engines.
- Add a reliability metric to each source API and use it to select the best source:
    - For example, a weighted-round-robin strategy.
- Add a "fastest" strategy that will gather response timing statistics and prefer the fastest provider.
- Support a choice of token bucket, fixed window, sliding window, leaky bucket, or burstable rate limiter algorithms.
- Consider supporting API providers purely by configuration files, rather than hardcoding them.
- Make max precision in results configurable. Currently defaults to 8 decimal places.
