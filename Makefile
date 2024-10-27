
# It is used during migration creations only
# like: "make migrate create create_users_table sql"
define set-goose-env
  GOOSE_DRIVER=DRIVER=sqlite3 \
	GOOSE_MIGRATION_DIR=./db/migrations \
	GOOSE_DBSTRING=./db/data/$(or $(1),$(APP_ENV)).db
endef

# Define a reusable command for loading .env variables
build:
	go build -o dist/server cmd/api/main.go
run:
	go run cmd/api/main.go
test:
	APP_ENV=test \
	go test -v -race -buildvcs ./... -count=1
migrate:
	$(set-goose-env) \
	go run cmd/migrate/main.go $(filter-out $@,$(MAKECMDGOALS))
migrate-test:
	APP_ENV=test \
	$(call set-goose-env,test) \
	go run cmd/migrate/main.go $(filter-out $@,$(MAKECMDGOALS))
migrate-build:
	go build -o dist/migrate cmd/migrate/main.go
migrate-run:
	$(call set-goose-env) \
	./dist/migrate $(filter-out $@,$(MAKECMDGOALS))
vet:
	go vet ./...


# Prevents Make from treating the arguments as targets
%:
	@:
