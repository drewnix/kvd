
build:
	go fmt .
	go build -o kvd main.go

run:
	go run main.go

test:
	go test -v ./...

test-cover:
	go test -cover -v ./...
	go test -coverprofile=coverage.out
	go tool cover -html coverage.out

clean:
	rm kvd
	rm coverage.out
