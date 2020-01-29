GO=go
BIN_NAME=mu2
MAIN_FILE=cmd/mu2/main.go

all: build

build:
	$(GO) build -o $(BIN_NAME) $(MAIN_FILE)

run: build
	./$(BIN_NAME)

docker-build-develop:
	docker build -t mu2:develop .

docker-run-develop:
	docker run mu2:develop

sloc:
	wc -l $$(find . -name '*.go')