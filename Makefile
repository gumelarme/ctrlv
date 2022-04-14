GOCMD=GOOS=linux CGO_ENABLED=0 go
GOBUILD= $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
    
all: clean build package
build: 
	$(GOBUILD) -o bin/$(BINARY_NAME) .
package:
	cp config.toml bin/
	cp -R public bin/
	cd bin && zip -r main.zip $(BINARY_NAME) config.toml public/
test: 
	$(GOTEST) -v ./...
clean: 
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_NAME).zip
