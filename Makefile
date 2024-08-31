# test output directory for test artifacts
test-output-dir = .test

# run all targets
all: init test-setup test

# target to initialize the directories needed for other targets ensuring linter is installed for lint directive
init:
	mkdir $(test-output-dir)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	golangci-lint run --config .\.golangci.yaml -v

# target to setup a devcontainer and a localstack on a docker network
dev-setup:
	docker compose up --wait

# target to setup a localstack on a docker network
# 	creates 3rd party stores for testing, IE: localstack as defined in the docker-compose file
test-setup:   
	docker compose up -d localstack --wait

# target to run test suite
#	run tests
#	generate new test artifacts in the test-output-dir
test:
	go test -cover -coverprofile=$(test-output-dir)/coverage.out -v ./... > $(test-output-dir)/test-run.log
	go tool cover -html=$(test-output-dir)/coverage.out -o $(test-output-dir)/coverage.html
