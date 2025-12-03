BINARY_NAME=argon
GO_BUILD=go build -o $(BINARY_NAME)
build:
        $(GO_BUILD)
install: build
        cp $(BINARY_NAME) /usr/local/bin/
clean:
        rm -f $(BINARY_NAME)
