.PHONY: openapi
openapi:
	oapi-codegen -generate types -o internal/categories/ports/openapi_types.gen.go -package ports api/openapi/categories.yaml
	oapi-codegen -generate server -o internal/categories/ports/openapi_api.gen.go -package ports api/openapi/categories.yaml
	oapi-codegen -generate spec -o internal/categories/ports/openapi_spec.gen.go -package ports api/openapi/categories.yaml

.PHONY: mocks
mocks:
	mockery --all --output testing/mocks  

.PHONY: lint
lint:
	golangci-lint run

