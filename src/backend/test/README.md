# Test Structure

Tests được tổ chức theo layer trong folder `test/`:

```
test/
├── service/
│   ├── auth/          # Unit tests cho auth service
│   ├── conversation/ # Unit tests cho conversation service
│   └── message/       # Unit tests cho message service
├── infra/
│   └── repo/         # Unit tests cho repository layer
├── integration/      # Integration tests (end-to-end)
│   ├── auth_test.go
│   ├── conversation_test.go
│   ├── message_test.go
│   └── helper.go
├── mocks/
│   └── repository.go # Shared mock repositories
├── endpoint/         # Tests cho endpoint layer
└── transport/
    └── http/         # Tests cho HTTP handlers
```

## Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific layer tests
make test-auth          # Auth service unit tests
make test-service       # All service unit tests
make test-repo          # Repository unit tests
make test-integration   # Integration tests (end-to-end)
```

## Test Types

### Unit Tests
- Test individual functions/methods in isolation
- Use mocks for dependencies
- Located in `test/service/`, `test/infra/repo/`, etc.

### Integration Tests
- Test complete flows from HTTP handler to database
- Use in-memory SQLite database
- Test real interactions between layers
- Located in `test/integration/`

## Test Helpers

- Mock repositories: `test/mocks/repository.go`
- Test helpers: `service/auth/test_helper.go`, `infra/repo/test_helper.go`
- Integration test setup: `test/integration/helper.go`

