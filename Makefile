CC=go
FLAGS=build

all: krawler

krawler:
	$(CC) $(FLAGS)

clean:
	rm krawler; rm /usr/local/bin/krawler

clean-bin:
	rm bin/*

install:
	cp krawler /usr/local/bin

compile-all:
	GOOS=darwin  GOARCH=amd64 go build -o bin/krawler-macos-amd64
	GOOS=darwin  GOARCH=arm64 go build -o bin/krawler-macos-arm64
	GOOS=linux   GOARCH=amd64 go build -o bin/krawler-linux-amd64
	GOOS=linux   GOARCH=386   go build -o bin/krawler-linux-i386
	GOOS=linux   GOARCH=arm64 go build -o bin/krawler-linux-arm64
	GOOS=windows GOARCH=amd64 go build -o bin/krawler-windows-amd64
	GOOS=windows GOARCH=386   go build -o bin/krawler-windows-i386