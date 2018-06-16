binary=wow-addon-downloader

SOURCE=$(shell find . -type f -name "*.go")

WINDOWS_BINARY=$(binary).exe
DARWIN_BINARY=$(binary)
TESTS=.tests

all: $(TESTS) $(WINDOWS_BINARY) $(DARWIN_BINARY)

.PHONY: windows
windows: $(WINDOWS_BINARY)

.PHONY: darwin
darwin: $(DARWIN_BINARY)

.PHONY: test
test: $(TESTS)

$(WINDOWS_BINARY): $(SOURCE)
	GOOS=windows GOARCH=amd64 go build -o $(binary).exe .

$(DARWIN_BINARY): $(SOURCE)
	GOOS=darwin GOARCH=amd64 go build -o $(binary) .

$(TESTS): $(SOURCE)
	go test ./...
	echo $$(date) >> $(TESTS)

.PHONY: clean
clean:
	rm -rf $(WINDOWS_BINARY) $(DARWIN_BINARY)  $(TESTS) vendor
