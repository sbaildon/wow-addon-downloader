binary=wow-addon-downloader

all: windows darwin

windows:
	GOOS=windows GOARCH=amd64 go build -o $(binary).exe .

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(binary) .