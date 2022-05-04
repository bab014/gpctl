.PHONY: build

build:
	go build -o gpctl && \
	GOARCH=amd64 GOOS=windows go build -o gpctl.exe
