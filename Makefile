all: clean reader writer piper

clean:
	go clean
	rm -rf bin/*

reader:
	go build -o bin/ ioTools/reader

writer:
	go build -o bin/ ioTools/writer

piper:
	go build -o bin/ ioTools/piper
