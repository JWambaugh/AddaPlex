PLUGINS := $(wildcard plugins/*)
GOCODE := $(wildcard *.go)


all: build plugins

build: bin/AddaPlex
	

bin/AddaPlex: bin gocode
	go build -o bin/AddaPlex

run: all
	./bin/AddaPlex


bin: 
	mkdir bin

gocode: $(GOCODE)

$(GOCODE):
	@echo $@

plugins: $(PLUGINS)
$(PLUGINS): bin
	go build -buildmode=plugin -o bin/$@.so $@/* 

clean:
	rm -rf bin