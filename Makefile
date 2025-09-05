# Build the application
.PHONY: build
build:
	@mkdir -p bin
	go build -o bin/mono-track ./src

# Run unit tests only (excluding e2e tests)
.PHONY: test-unit
test-unit:
	@mkdir -p coverage
	go test -v -json ./src > coverage/unit-tests.json
	@go run ./scripts/gotestdox-json/main.go coverage/unit-tests.json

# Run unit tests with coverage (excluding e2e tests)
.PHONY: coverage-unit
coverage-unit:
	@mkdir -p coverage
	@echo "Running unit tests with coverage..."
	go test -v -json -coverprofile=coverage/unit.out ./src > coverage/unit-tests.json
	@go run ./scripts/gotestdox-json/main.go coverage/unit-tests.json
	@echo ""
	@echo "Coverage Report by File:"
	@go run ./scripts/format-coverage/main.go coverage/unit.out

# Build with coverage instrumentation for e2e tests
.PHONY: build-e2e-coverage
build-e2e-coverage:
	@mkdir -p bin
	go build -cover -covermode=set -coverpkg=./... -o bin/mono-track-e2e ./src

# Run e2e tests
.PHONY: test-e2e
test-e2e: build-e2e-coverage
	@rm -rf coverage/e2e && mkdir -p coverage/e2e
	@echo "Running e2e tests..."
	@set -e; \
	if go test -v -json -timeout 5m ./tests/e2e/... > coverage/e2e-tests.json 2>&1; then \
		test_exit_code=0; \
		echo "E2E tests passed"; \
	else \
		test_exit_code=$$?; \
		echo "Tests failed with exit code $$test_exit_code, but continuing to show results..."; \
	fi; \
	echo ""; \
	echo "E2E Test Results:"; \
	go run ./scripts/gotestdox-json/main.go coverage/e2e-tests.json; \
	exit $$test_exit_code

# Generate e2e coverage report
.PHONY: coverage-e2e-report
coverage-e2e-report:
	@if [ -d coverage/e2e ] && [ "$$(find coverage/e2e -name 'covcounters*' -o -name 'covmeta*' | wc -l)" -gt 0 ]; then \
		rm -rf coverage/e2e-merged && mkdir -p coverage/e2e-merged && \
		covdirs="$$(find coverage/e2e -path 'coverage/e2e-merged' -prune -o -type f \( -name 'covmeta*' -o -name 'covcounters*' \) -exec dirname {} \; | sort -u | tr '\n' ',' | sed 's/,$$//')"; \
		go tool covdata merge -i="$$covdirs" -o=coverage/e2e-merged && \
		go tool covdata textfmt -i=coverage/e2e-merged -o=coverage/e2e.out; \
	else \
		echo "No e2e coverage data found"; \
	fi

# E2E coverage (runs tests and generates reports)
.PHONY: coverage-e2e
coverage-e2e: build-e2e-coverage
	@mkdir -p coverage/e2e
	@echo "Running e2e tests with coverage..."
	@set -e; \
	if go test -v -json -timeout 5m ./tests/e2e/... > coverage/e2e-tests.json 2>&1; then \
		test_exit_code=0; \
		echo "E2E tests with coverage passed"; \
	else \
		test_exit_code=$$?; \
		echo "Tests failed with exit code $$test_exit_code, but continuing to show results..."; \
	fi; \
	$(MAKE) coverage-e2e-report; \
	echo ""; \
	echo "E2E Test Results:"; \
	go run ./scripts/gotestdox-json/main.go coverage/e2e-tests.json; \
	echo ""; \
	echo "E2E Coverage Report by File:"; \
	if [ -f coverage/e2e.out ]; then \
		go run ./scripts/format-coverage/main.go coverage/e2e.out; \
	else \
		echo "No e2e coverage data available"; \
	fi; \
	exit $$test_exit_code

# Get uncovered lines from SonarCloud for current PR
.PHONY: sonar-uncovered
sonar-uncovered:
	go run ./scripts/sonar-uncovered-lines/main.go

# Clean Git branches and project state (reverts all changes, removes untracked files, deletes all branches except main and current)
.PHONY: clean-branches
clean-branches:
	@echo "Cleaning up Git branches and project..."
	@current_branch=$$(git branch --show-current); \
	echo "Current branch: $$current_branch"; \
	echo "Reverting all changes to tracked files..."; \
	git reset --hard HEAD; \
	echo "Removing untracked files and directories..."; \
	git clean -fd; \
	branches_to_delete=$$(git branch --format='%(refname:short)' | grep -v "^main$$" | grep -v "^$$current_branch$$"); \
	if [ -n "$$branches_to_delete" ]; then \
		echo "Branches to delete:"; \
		echo "$$branches_to_delete"; \
		echo "$$branches_to_delete" | xargs -r git branch -D; \
		echo "Branch cleanup complete"; \
	else \
		echo "No branches to delete"; \
	fi; \
	echo "Project cleanup complete"

# Clean build files
.PHONY: clean
clean:
	rm -rf bin/ coverage/
