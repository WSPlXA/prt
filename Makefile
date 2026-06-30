build:
	go build -o ./bin/prt

test:
	go test -v ./...

tidy:
	go mod tidy