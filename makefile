

build: build-daemon build-client

build-daemon:
	go build -v -o bin/vaultd .

build-client:
	go build -v -o bin/vault ./cmd

clean:
	@rm -dRf bin