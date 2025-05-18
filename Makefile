build:
	go build -C cmd/gatherer/ -o ../../bin/

run: build
	./bin/gatherer -env ./.env
