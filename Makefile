build:
	go build -o dist/server cmd/api/main.go
run:
	go run cmd/api/main.go
test:
	go test -v ./... -count=1
migrate:
	go run cmd/migrate/main.go $(GOOSE_DRIVER) $(DB_URL) $(filter-out $@,$(MAKECMDGOALS))
migrate-build:
	go build -o dist/migrate cmd/migrate/main.go
migrate-run:
	./dist/migrate $(GOOSE_DRIVER) $(DB_URL) $(filter-out $@,$(MAKECMDGOALS))


# Prevents Make from treating the arguments as targets
%:
	@:
