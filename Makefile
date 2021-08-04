default: 

run-server:

run-import: build
	./build/import

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o build/our-expenses-server ./cmd/server
	CGO_ENABLED=0 GOARCH=amd64 go build -o build/import ./cmd/import

test:
	go test ./... -cover

openapi:
	oapi-codegen -generate types -o internal/categories/ports/openapi_types.gen.go -package ports api/openapi/categories.yaml
	oapi-codegen -generate server -o internal/categories/ports/openapi_api.gen.go -package ports api/openapi/categories.yaml
	oapi-codegen -generate spec -o internal/categories/ports/openapi_spec.gen.go -package ports api/openapi/categories.yaml

	oapi-codegen -generate types -o internal/expenses/ports/openapi_types.gen.go -package ports api/openapi/expenses.yaml
	oapi-codegen -generate server -o internal/expenses/ports/openapi_api.gen.go -package ports api/openapi/expenses.yaml
	oapi-codegen -generate spec -o internal/expenses/ports/openapi_spec.gen.go -package ports api/openapi/expenses.yaml

mocks:
	mockery --all --output testing/mocks  

lint:
	golangci-lint run

