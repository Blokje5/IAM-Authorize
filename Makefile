# BUILD
build: bin
	go build -o bin/iam
bin:
	mkdir -p bin

# DEVELOPMENT ENVIRONMENT
result = $(shell docker ps --format '{{.Names}}' | grep iam-postgres)
start_local_db:
ifeq ($(result),iam-postgres)
	echo "Already running"
else
	docker run --name iam-postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=iam_test -e POSTGRES_PASSWORD=local -p 5432:5432 -d postgres:12
endif

PG_URL ?= "postgresql://postgres@127.0.0.1:5432/iam_test?sslmode=disable&password=local"
run_local: start_local_db
	go run main.go run --connectionstring $(PG_URL)

clean_local:
	docker rm -f iam-postgres || true

# TESTING
test: start_local_db
	export TEST_CONNECTIONSTRING=$(PG_URL) && go test ./...

.PHONY: build start_local_db run_local clean_local test result