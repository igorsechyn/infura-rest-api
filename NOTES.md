# Design

For service design I followed clean architecture principle to separate frameworks from business logic. Main components of the application are:

- `web_server.go` - handles incoming requests to retrieve information. Main responsibility is to extract request input (path parameters, query, etc.), call business logic unit and process response (map errors to status code, write payload). I chose `chi` web framework, as I was the most familiar with it. However it should be super simple to replace it with any other framework
- `app.go` - container that holds all business unit components, like `config`, `reporter`, various services for processing transactions, blocks, etc.
- `transactions` package contains a service module, which is responsible for callinig transactions client to retrieve information and reporting on the results (error, success). transactions clien is using an `ethereum client` to interact with ethereum. At runtime an infura imlementation of the client is injected. For unit tests a mock client is injected to test business logic
- `reporter` module is responsible for reporting different events. It may contain multiple sinks to report events to different targets. At runtime a json logger is initialised for reporting. But there can also be a metrics sink to report events to statsd for example
- `config` module provides application configuration. At runtime it is read from environment variables
- `ethereum` package defines an interface to interact with ethereum. `infura` provides an implementation using infura service 

# Testing

## Unit testing

I am huge fan of TDD and behavioural tests. Therefore my major goal is always to decouple tests from internal implementation details and code structure and make them sensitive only to behaviour changes. e.g. `transactions-client.go` module is not tested in isolation, I only extracted it after I wrote all the logic and tests inside the service module.

I am using dependency injection to inject boundaries, like http clients, fs clients, etc... and assert behaviour on those boundaries

## E2E tests

I am usually trying to keep the number of e2e tests as low as possible, they tend to be slow, flakey and not targeted. Major purpose is to verify the implementation of the boundaries, like http clients (in this case infura client). I am not using a real infura service for e2e tests, to avoid flakiness. Instead I start a small mock server. There is currently nothing implemented to verify the correctness of the mock, but I have ideas, we can talk about :) just did not have enough time to implement them.

## Performance tests

I used artillery to run performance, which introduced a dependency on node, but this is the tool I am most familiar with. I ran following scenario:

- ramp up from 10 to 30 concurrent users over 60 seconds
- sustain load of 30 concurrent users for 5 minutes
- Requests were alternating between 30 different block numbers

Results are under `./test/performance/2020-01-27---10-32-28`:

- max latency was 3.5K ms
- p99 latency was 636 ms

I was not able to run higher loads, because my laptop was not able to keep up. One of the reasons why I chose artillery, because it offers to run distributed tests on AWS Fahrgate. But I havent done that, as it is only available in the pro version.

I did not see any rate limits from infura, I assume the data I was requesting is served from cache