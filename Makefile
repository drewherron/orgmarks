VERSION ?= 0.1.0

.PHONY: release-all clean-dist

release-all: clean-dist
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/orgmarks-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/orgmarks-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o dist/orgmarks-macos-intel .
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o dist/orgmarks-macos-arm64 .

clean-dist:
	rm -rf dist/
