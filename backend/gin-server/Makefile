# make sure to use tabs instead of spaces for identation

run:
	cd server && nodemon

test:
	go test ./server/internal/...

build:
	rm -rf ./bin
	cd server && go mod tidy
	go build -o bin/server ./server/cmd/server
	cd server && cp .env.dev ../bin/
	cd server && cp flowing-castle-447017-h0-72697064ea74.json ../bin/

# OpenAPI targets
openapi-copy:
	@echo "Copying OpenAPI spec to server directory..."
	mkdir -p server/openapi
	cp ../../openapi/api.yaml server/openapi/

openapi-serve:
	@echo "Starting server with Swagger UI at http://localhost:8080/swagger/"
	cd server && OPENAPI_SPEC_PATH=./openapi/api.yaml go run cmd/server/main.go
