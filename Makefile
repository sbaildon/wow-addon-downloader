binary=wow-addon-downloader

SOURCE=$(shell find . -type f -name "*.go")

WINDOWS_BINARY=$(binary).exe
DARWIN_BINARY=$(binary)

all: $(WINDOWS_BINARY) $(DARWIN_BINARY)

.PHONY: windows
windows: $(WINDOWS_BINARY)

.PHONY: darwin
darwin: $(DARWIN_BINARY)

$(WINDOWS_BINARY): $(SOURCE)
	GOOS=windows GOARCH=amd64 go build -o $(binary).exe .

$(DARWIN_BINARY): $(SOURCE)
	GOOS=darwin GOARCH=amd64 go build -o $(binary) .

.PHONY: clean
clean:
	rm $(WINDOWS_BINARY) $(DARWIN_BINARY)