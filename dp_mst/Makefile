NAME=main
PWD=$( pwd )

build:
	mkdir -p bin
	GOARCH=amd64 GOPATH=${PWD} go build -o ./bin/${NAME} cmd/cmdmst/main.go

run: build
	echo 1 | ./bin/${NAME} --file=./samples/example

all: build

clean:
	go clean
	rm -rf ./bin