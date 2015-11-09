PREFIX := /usr/local

clientbin := vault
daemonbin := vaultd

build: build-daemon build-client

build-daemon:
	go build -v -o bin/$(daemonbin) .

build-client:
	go build -v -o bin/$(clientbin) ./cmd


install:
	#test -d $(PREFIX) || mkdir $(PREFIX)
  #test -d $(PREFIX)/bin || mkdir $(PREFIX)/bin
	install -m 0755 bin/$(clientbin) $(PREFIX)/bin
	install -m 0755 bin/$(daemonbin) $(PREFIX)/bin

uninstall:
	rm $(PREFIX)/bin/$(clientbin)
	rm $(PREFIX)/bin/$(daemonbin)

clean:
	@rm -dRf bin

