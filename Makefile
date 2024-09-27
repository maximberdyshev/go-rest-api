# run program without compile
dev:
	go run ./cmd/go-rest-api/main.go

# compile program
build:
	rm -rf ./build
	go build -o ./build/go-rest-api ./cmd/go-rest-api

# run compiled program
start:
	./build/go-rest-api

# clean compiled program
clean:
	rm -rf ./build

# Generate swagger docs
genDoc:
	swag init -d cmd/go-rest-api/,internal/transport/http/v1/handler/ -g main.go --parseDependency

# Format swagger comments
fmtDoc:
	swag fmt -d cmd/go-rest-api/,internal/transport/http/v1/handler/