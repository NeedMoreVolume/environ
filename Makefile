# test output directory for test artifacts
test-output-dir = .test

# run all targets
all: init test-setup test

# target to initialize the directories needed for other targets
init:
	mkdir $(test-output-dir)

# target to setup test environcment
# 	creates 3rd party stores for testing, IE: localstack as defined in the docker-compose file
test-setup:   
	docker compose up --wait

# target to run test suite
#	run tests
#	generate new test artifacts in the test-output-dir
test:
	go test -cover -coverprofile=$(test-output-dir)/coverage.out -v ./... > $(test-output-dir)/test-run.log
	go tool cover -html=$(test-output-dir)/coverage.out -o $(test-output-dir)/coverage.html
