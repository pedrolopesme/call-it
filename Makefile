GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD)fmt
BINARY_NAME=$(GOPATH)/bin/call-it
BINARY_UNIX=$(BINARY_NAME)_unix

build: 
	@echo "Building call-it"
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	@echo "Running call-it tests"
	$(GOTEST) -v ./...

clean: 
	@echo "Cleaning call-it"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	@echo "Running call-it with params: $(filter-out $@,$(MAKECMDGOALS))"
	$(GOBUILD) -o $(BINARY_NAME)
	./$(BINARY_NAME) $(filter-out $@,$(MAKECMDGOALS))

fmt:
	@echo "Running gofmt for all project files"
	$(GOFMT) -w *.go