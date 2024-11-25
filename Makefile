lint:
	golangci-lint run --fix ./...

migrate:
	docker run --rm -v ./migrations:/migrations  migrate/migrate:v4.14.1 -path=/migrations -database="postgres://myuser:mypassword@host.docker.internal:5432/mydatabase?sslmode=disable" up

build-binary:
	go build -o bin/main cmd/main.go
