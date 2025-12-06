.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -o ./docs --parseDependency --parseInternal --parseDepth 2
	@echo "Swagger documentation generated successfully!"
	@echo "Access Swagger UI at: http://localhost:80/swagger/index.html"

.PHONY: swagger-minimal
swagger-minimal:
	@echo "Generating minimal Swagger documentation..."
	@swag init -g main.go -o ./docs --parseDependency --parseInternal --parseDepth 1 --generatedTime false
	@echo "Minimal Swagger documentation generated!"

.PHONY: run
run:
	@go run main.go

.PHONY: build
build:
	@go build -o bin/server main.go

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./test/... -v

.PHONY: test-progress
test-progress:
	@go run scripts/test_progress.go

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./test/... -cover

.PHONY: test-auth
test-auth:
	@echo "Running auth service tests..."
	@go test ./test/service/auth/... -v

.PHONY: test-service
test-service:
	@echo "Running all service tests..."
	@go test ./test/service/... -v

.PHONY: test-repo
test-repo:
	@echo "Running repository tests..."
	@go test ./test/infra/repo/... -v

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@go test ./test/integration/... -v

