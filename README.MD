## Run and test

Run:
```bash
   cd deployment
   make run
   cd ../
```

Migrate:
```bash
   make migrate
```

Wait for all services to start. docker ps should show all containers running.

## Solution Design

- I used https://en.wikipedia.org/wiki/Hexagonal_architecture_(software) architecture for the solution.
- Ideas for microservice template are taken from https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example repository.

## Dependencies

- grpc dependencies are used for communication between services.
- opentelemetry, prometheus for metrics and tracing. // I din't implement it because it is not required in the task.
- github.com/lib/pq  for working with postgres.
- github.com/oklog/run for implementing graceful shutdown easily.