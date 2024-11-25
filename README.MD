## Run and test

Build:
```bash
  cd deployment
  make build
  cd ../
```

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

## How to test
Create accounts
```bash
grpcurl -vv -plaintext -d '{
"id": "36C980FB-F1FE-42F4-846A-ADB134027D59",
"owner_id": "46C980FB-F1FE-42F4-846A-ADB134027D50",
"currency": 1
}' localhost:50055 accountservice.v1beta1.AccountService.CreateAccount
```

```bash
grpcurl -vv -plaintext -d '{
"id": "46C980FB-F1FE-42F4-846A-ADB134027D50",
"owner_id": "46C980FB-F1FE-42F4-846A-ADB134027D50",
"currency": 1
}' localhost:50055 accountservice.v1beta1.AccountService.CreateAccount
```

Top up an account
```bash
grpcurl -vv -plaintext -d '{
"id": "36C980FB-F1FE-42F4-846A-ADB134027D59",
"amount": 100
}' localhost:50055 accountservice.v1beta1.AccountService.TopUpAccount
```

Transfer:
```bash
grpcurl -vv -plaintext -d '{
"from_account_id": "36C980FB-F1FE-42F4-846A-ADB134027D59",
"to_account_id": "46C980FB-F1FE-42F4-846A-ADB134027D50",
"amount": 43
}' localhost:50055 accountservice.v1beta1.AccountService.TransferFunds
```

Get Accounts:
```bash
grpcurl -vv -plaintext -d '{
"id": "36C980FB-F1FE-42F4-846A-ADB134027D59"
}' localhost:50055 accountservice.v1beta1.AccountService.GetAccount

```

```bash
grpcurl -vv -plaintext -d '{
"id": "46C980FB-F1FE-42F4-846A-ADB134027D50"
}' localhost:50055 accountservice.v1beta1.AccountService.GetAccount

```

## Solution Design

- I used https://en.wikipedia.org/wiki/Hexagonal_architecture_(software) architecture for the solution.
- Ideas for microservice template are taken from https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example repository.

I already have a template for CQRS and DDD and to be able to implement the solution in short time, I used my ready template.
I used GRPC for API because again I already have a template for it and it was easy to implement. Also, if this API for internal use, GRPC is a good choice.

- [adapters](internal/adapters) - contains the implementation of the interfaces defined in the application, domain packages. It is responsible for the communication with the external world.
- [ports](internal/ports) - contains implementation how to communicate with the service.
- [domain](internal/domain) - business logic of the service.
- [app](internal/app) - application layer, contains the implementation of the use cases.
- [common](internal/common) - contains common code, that can be reused in different microservices.
- [service](internal/service) - contains setup for microservice and DI.
- [component_test.go](internal/service/component_test.go) - contains component tests for the service.

[deployment](deployment) - contains docker-compose file and Makefile for building and running the service.
[api](api) - contains proto files for the service.
[server](cmd/server) - contains the main function for the service.

## Dependencies

- grpc dependencies are used for communication between services.
- opentelemetry, prometheus for metrics and tracing. // I din't implement it because it is not required in the task.
- github.com/lib/pq  for working with postgres.
- github.com/oklog/run for implementing graceful shutdown easily.