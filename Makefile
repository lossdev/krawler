CC=go
FLAGS=build

all: krawler

krawler:
	$(CC) $(FLAGS)

clean:
	rm krawler; rm /usr/local/bin/krawler

install:
	cp krawler /usr/local/bin