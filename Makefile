build:
	go build -o ./bin/treechat

run: build
	./bin/treechat

test:
	go test -o ./...