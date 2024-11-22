
# Define a reusable command for loading .env variables
build:
	go build -o dist/server cmd/api/main.go
	go build -o dist/migrate cmd/migrate/main.go
run:
	go run cmd/api/main.go
test:
	APP_ENV=test \
	go test -v -race -buildvcs ./... -count=1
db-schema-dump:
	go run cmd/db_schema_dump/main.go
db-schema-check:
	go run cmd/db_schema_check/main.go
migrate:
	go run cmd/migrate/main.go $(filter-out $@,$(MAKECMDGOALS))
migrate-test:
	APP_ENV=test \
	go run cmd/migrate/main.go $(filter-out $@,$(MAKECMDGOALS))
migrate-run:
	./dist/migrate $(filter-out $@,$(MAKECMDGOALS))
vet:
	go vet ./...


# Prevents Make from treating the arguments as targets
%:
	@:
