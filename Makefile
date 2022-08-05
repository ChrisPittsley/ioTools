all: clean reader writer piper

clean:
	go clean
	rm -rf bin/*

reader:
	go build -buildvcs=false -o bin/ ioTools/reader

writer:
	go build -buildvcs=false -o bin/ ioTools/writer

piper:
	go build -buildvcs=false -o bin/ ioTools/piper
