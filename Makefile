.PHONY: build run create-mock create-dir

build:
	go build -o bin/file-server cmd/web/main.go && cd website && npm run build

run:
	./bin/file-server

create-mock:
	echo "Hello World" > shared/test.txt

create-dir:
	mkdir shared && chmod 755 shared