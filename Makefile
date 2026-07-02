build:
	go build -o ./bin/prt

run build:
	./bin/prt
test:
	go test -v ./...

tidy:
	go mod tidy