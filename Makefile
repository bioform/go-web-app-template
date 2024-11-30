# Change these variables as necessary.
main_package_path = ./cmd/api/main.go
binary_name = server

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test vet lint/all
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"

## test: run all tests
.PHONY: test
test:
	@ginkgo -r --randomize-all --randomize-suites -cover -coverprofile=coverage.out --output-dir=.

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover: test
	@go tool cover -html=coverage.out -o coverage.html


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build: build/migration/up
	go build -o dist/${binary_name} ${main_package_path}

## build: build the application
.PHONY: build/migration/up
build/migration/up:
	go build -o dist/migration-up cmd/migration_up/main.go

## run: run the API server
.PHONY: run
run: build
	./dist/server

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "dist/${binary_name}" --build.delay "100" \
		--build.exclude_dir "db,dist" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico, out" \
		--misc.clean_on_exit "true"

## db/schema/dump: dump the database schema
.PHONY: db/schema/dump
db/schema/dump:
	go run cmd/db_schema_dump/main.go

## db/schema/check: check the database schema
db/schema/check:
	go run cmd/db_schema_check/main.go

## db/migration: run database migrations. Usage: make db/migration up/down
.PHONY: db/migration
db/migration:
	go run cmd/migration/main.go $(filter-out $@,$(MAKECMDGOALS))

## db/test/prepare: prepare the test database. Recreate the test database and run restore from the schema dump
.PHONY: db/test/prepare
db/test/prepare:
	APP_ENV=test \
	go run cmd/db_test_prepare/main.go

## db/migration/up: run database migrations from the binary
.PHONY: db/migration/up
db/migration/up: build/migration/up
	dist/migration-up

## vet: run go vet
.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
## lint: runs linter for a given directory
lint:
	@ if [ -z "$(PACKAGE)" ]; then echo >&2 please set directory via variable PACKAGE; exit 2; fi
	@ docker run -t  --rm -v "`pwd`:/workspace:cached" -w "/workspace/$(PACKAGE)" golangci/golangci-lint:latest golangci-lint run


.PHONY: lint/all
## lint/all: runs linter for all packages
lint/all:
	@ docker run -t  --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:latest golangci-lint run

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
	upx -5 /tmp/bin/linux_amd64/${binary_name}
	# Include additional deployment steps here...

# Prevents Make from treating the arguments as targets
# Allows to run targets with arguments
TARGETS_WITH_ARGS := db/migration

%:
	@if ! echo "$(TARGETS_WITH_ARGS)" | grep -qw "$(word 1, $(MAKECMDGOALS))"; then \
		echo "Error: Unknown target '$@'"; \
		exit 1; \
	fi
	@:
