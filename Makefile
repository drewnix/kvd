
build:
	go fmt ./...
	go build -o kv main.go

run:
	go run main.go

serve:
	go run main.go serve

test:
	go test -v ./...

test-cover:
	go test -cover -v ./...
	go test -coverprofile=coverage.out
	go tool cover -html coverage.out

clean:
	rm kv
	rm coverage.out
