lint:
	golangci-lint run --fix ./...

migrate:
	docker run --rm -v ./migrations:/migrations  migrate/migrate:v4.14.1 -path=/migrations -database="postgres://myuser:mypassword@host.docker.internal:5432/mydatabase?sslmode=disable" up

build-binary:
	go build -o bin/main cmd/main.go

image:
	docker build -t gcr.io/dsxack/accounts:latest .
	docker push gcr.io/dsxack/accounts:latest


ssh:
	@ssh dmitriysmotrov@34.32.31.67

deploy:
	cd deployment && docker compose up -d