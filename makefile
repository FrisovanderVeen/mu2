GO=go
BIN_NAME=mu2

all: build

build:
	$(GO) build -o $(BIN_NAME) main.go

run: build
	./$(BIN_NAME)

sloc:
	wc -l $$(find . -name '*.go')